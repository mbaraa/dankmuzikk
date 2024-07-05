package pages

import (
	"context"
	"dankmuzikk/config"
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/handlers/middlewares/auth"
	"dankmuzikk/handlers/middlewares/contenttype"
	"dankmuzikk/log"
	"dankmuzikk/models"
	"dankmuzikk/services/history"
	"dankmuzikk/services/jwt"
	"dankmuzikk/services/playlists"
	"dankmuzikk/services/playlists/songs"
	"dankmuzikk/services/youtube/download"
	"dankmuzikk/services/youtube/search"
	"dankmuzikk/views/components/status"
	"dankmuzikk/views/layouts"
	"dankmuzikk/views/pages"
	"errors"
	"net/http"

	_ "github.com/a-h/templ"
)

type pagesHandler struct {
	profileRepo      db.GetterRepo[models.Profile]
	playlistsService *playlists.Service
	jwtUtil          jwt.Manager[jwt.Json]
	ytSearch         search.Service
	downloadService  *download.Service
	historyService   *history.Service
	songsService     *songs.Service
}

func NewPagesHandler(
	profileRepo db.GetterRepo[models.Profile],
	playlistsService *playlists.Service,
	jwtUtil jwt.Manager[jwt.Json],
	ytSearch search.Service,
	downloadService *download.Service,
	historyService *history.Service,
	songsService *songs.Service,
) *pagesHandler {
	return &pagesHandler{
		profileRepo:      profileRepo,
		playlistsService: playlistsService,
		jwtUtil:          jwtUtil,
		ytSearch:         ytSearch,
		downloadService:  downloadService,
		historyService:   historyService,
		songsService:     songsService,
	}
}

func (p *pagesHandler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	var recentPlays []entities.Song
	var err error
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if profileIdCorrect {
		recentPlays, err = p.historyService.Get(profileId, 1)
		if err != nil {
			log.Errorln(err)
		}
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Home")
		w.Header().Set("HX-Push-Url", "/")
		pages.Index(recentPlays).Render(r.Context(), w)
		return
	}
	layouts.Default("Home", pages.Index(recentPlays)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "About")
		w.Header().Set("HX-Push-Url", "/about")
		pages.About().Render(r.Context(), w)
		return
	}
	layouts.Default("About", pages.About()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw("Login", pages.Login()).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		status.
			BugsBunnyError("I'm not sure what you're trying to do :)").
			Render(context.Background(), w)
		return
	}

	playlists, err := p.playlistsService.GetAll(profileId)
	if err != nil {
		playlists = make([]entities.Playlist, 0)
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Playlists")
		w.Header().Set("HX-Push-Url", "/playlists")
		pages.Playlists(playlists).Render(r.Context(), w)
		return
	}
	layouts.Default("Playlists", pages.Playlists(playlists)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSinglePlaylistPage(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		status.
			BugsBunnyError("I'm not sure what you're trying to do :)").
			Render(context.Background(), w)
		return
	}

	playlistPubId := r.PathValue("playlist_id")
	if playlistPubId == "" {
		status.
			BugsBunnyError("You need to provide a playlist id!").
			Render(context.Background(), w)
		return
	}

	playlist, permission, err := p.playlistsService.Get(playlistPubId, profileId)
	htmxReq := contenttype.IsNoLayoutPage(r)
	switch {
	case errors.Is(err, playlists.ErrUnauthorizedToSeePlaylist):
		log.Errorln(err)
		if htmxReq {
			status.
				BugsBunnyError("You can't see this playlist! <br/> (don't snoop around other people's stuff or else!)").
				Render(context.Background(), w)
		} else {
			layouts.Default("Error",
				status.
					BugsBunnyError("You can't see this playlist! <br/> (don't snoop around other people's stuff or else!)")).
				Render(r.Context(), w)
		}
		return
	case err != nil:
		if htmxReq {
			status.
				BugsBunnyError("You can't see this playlist! <br/> (it might be John Cena)").
				Render(context.Background(), w)
		} else {
			layouts.Default("Error",
				status.
					BugsBunnyError("You can't see this playlist! <br/> (it might be John Cena)")).
				Render(r.Context(), w)
		}
	}
	ctx := context.WithValue(r.Context(), auth.PlaylistPermission, permission)

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", playlist.Title)
		w.Header().Set("HX-Push-Url", "/playlist/"+playlist.PublicId)
		pages.Playlist(playlist).Render(ctx, w)
		return
	}
	layouts.Default(playlist.Title, pages.Playlist(playlist)).Render(ctx, w)
}

func (p *pagesHandler) HandleSingleSongPage(w http.ResponseWriter, r *http.Request) {
	songId := r.PathValue("song_id")
	if songId == "" {
		status.
			BugsBunnyError("You need to provide a song id!").
			Render(context.Background(), w)
		return
	}

	song, err := p.songsService.GetSong(songId)
	if err != nil {
		status.
			BugsBunnyError("Song doesn't exist!").
			Render(context.Background(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", song.Title)
		w.Header().Set("HX-Push-Url", "/song/"+song.YtId)
		pages.Song(song).Render(r.Context(), w)
		return
	}
	layouts.Default(song.Title, pages.Song(song)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	layouts.Default("Privacy", pages.Privacy()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	profileId, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		if contenttype.IsNoLayoutPage(r) {
			w.Header().Set("HX-Redirect", "/")
		} else {
			http.Redirect(w, r, config.Env().Hostname, http.StatusTemporaryRedirect)
		}
		return
	}
	// error is ignored, because the id was checked in the AuthHandler
	dbProfile, _ := p.profileRepo.Get(profileId)
	profile := entities.Profile{
		Name:     dbProfile.Name,
		PfpLink:  dbProfile.PfpLink,
		Username: dbProfile.Username,
	}
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Profile")
		w.Header().Set("HX-Push-Url", "/profile")
		pages.Profile(profile).Render(r.Context(), w)
		return
	}
	layouts.Default("Profile", pages.Profile(profile)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSearchResultsPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results, err := p.ytSearch.Search(query)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		log.Errorln(err)
		return
	}

	if len(results) != 0 {
		// TODO: move this call out of here
		log.Info("downloading songs' meta data from search")
		_ = p.downloadService.DownloadYoutubeSongsMetadata(results)
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Results for "+query)
		w.Header().Set("HX-Push-Url", "/search?query="+query)
		pages.SearchResults(results).Render(r.Context(), w)
		return
	}
	layouts.Default("Results for "+query, pages.SearchResults(results)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw("Sign up", pages.Signup()).Render(r.Context(), w)
}
