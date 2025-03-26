package pages

import (
	"context"
	"dankmuzikk-web/config"
	"dankmuzikk-web/entities"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/handlers/middlewares/contenttype"
	"dankmuzikk-web/log"
	"dankmuzikk-web/services/history"
	"dankmuzikk-web/services/playlists"
	"dankmuzikk-web/services/playlists/songs"
	"dankmuzikk-web/services/requests"
	"dankmuzikk-web/services/youtube/search"
	"dankmuzikk-web/views/components/status"
	"dankmuzikk-web/views/layouts"
	"dankmuzikk-web/views/pages"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/a-h/templ"
)

type pagesHandler struct {
	playlistsService *playlists.Service
	ytSearch         search.Service
	historyService   *history.Service
	songsService     *songs.Service
}

func NewPagesHandler(
	playlistsService *playlists.Service,
	ytSearch search.Service,
	historyService *history.Service,
	songsService *songs.Service,
) *pagesHandler {
	return &pagesHandler{
		playlistsService: playlistsService,
		ytSearch:         ytSearch,
		historyService:   historyService,
		songsService:     songsService,
	}
}

func (p *pagesHandler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	var recentPlays []entities.Song
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if profileIdCorrect {
		sessionToken, err := r.Cookie(auth.SessionTokenKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		recentPlays, err = p.historyService.Get(sessionToken.Value, 1)
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
	layouts.Default(layouts.PageProps{
		Title:       "Home",
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname,
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Index(recentPlays)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleAboutPage(w http.ResponseWriter, r *http.Request) {
	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "About")
		w.Header().Set("HX-Push-Url", "/about")
		pages.About().Render(r.Context(), w)
		return
	}
	layouts.Default(layouts.PageProps{
		Title:       "About",
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/about",
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.About()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw(layouts.PageProps{
		Title:       "Login",
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/login",
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Login()).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		status.
			BugsBunnyError("I'm not sure what you're trying to do :)").
			Render(context.Background(), w)
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlists, err := p.playlistsService.GetAll(sessionToken.Value)
	if err != nil {
		log.Errorln(err)
		playlists = make([]entities.Playlist, 0)
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Playlists")
		w.Header().Set("HX-Push-Url", "/playlists")
		pages.Playlists(playlists).Render(r.Context(), w)
		return
	}
	layouts.Default(layouts.PageProps{
		Title:       "Playlists",
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/playlists",
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Playlists(playlists)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSinglePlaylistPage(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		status.
			BugsBunnyError("I'm not sure what you're trying to do :)").
			Render(context.Background(), w)
		return
	}
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	playlistPubId := r.PathValue("playlist_id")
	if playlistPubId == "" {
		status.
			BugsBunnyError("You need to provide a playlist id!").
			Render(context.Background(), w)
		return
	}

	playlist, err := p.playlistsService.Get(sessionToken.Value, playlistPubId)
	htmxReq := contenttype.IsNoLayoutPage(r)
	switch {
	case errors.Is(err, playlists.ErrUnauthorizedToSeePlaylist):
		log.Errorln(err)
		if htmxReq {
			status.
				BugsBunnyError("You can't see this playlist! <br/> (don't snoop around other people's stuff or else!)").
				Render(context.Background(), w)
		} else {
			layouts.Default(layouts.PageProps{
				Title: "Error",
			},
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
			layouts.Default(layouts.PageProps{
				Title: "Error",
			},
				status.
					BugsBunnyError("You can't see this playlist! <br/> (it might be John Cena)")).
				Render(r.Context(), w)
		}
	}
	ctx := context.WithValue(r.Context(), auth.PlaylistPermission, playlist.Permissions)

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", playlist.Title)
		w.Header().Set("HX-Push-Url", "/playlist/"+playlist.PublicId)
		pages.Playlist(playlist).Render(ctx, w)
		return
	}
	layouts.Default(layouts.PageProps{
		Title:       playlist.Title,
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/playlist/" + playlist.PublicId,
		Type:        layouts.PlaylistPage,
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Playlist(playlist)).Render(ctx, w)
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
	layouts.Default(layouts.PageProps{
		Title:       song.Title,
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/song/" + song.YtId,
		Type:        layouts.SongPage,
		ImageUrl:    song.ThumbnailUrl,
		Audio: layouts.AudioProps{
			Url:      fmt.Sprintf("%s/muzikkx/%s.mp3", config.Env().CdnAddress, song.YtId),
			Duration: song.Duration,
			Musician: song.Artist,
		},
	}, pages.Song(song)).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	layouts.Default(layouts.PageProps{
		Title:       "Privacy",
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/privacy",
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Privacy()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	_, profileIdCorrect := r.Context().Value(auth.ProfileIdKey).(uint)
	if !profileIdCorrect {
		if contenttype.IsNoLayoutPage(r) {
			w.Header().Set("HX-Redirect", "/")
		} else {
			http.Redirect(w, r, "https://"+config.Env().Hostname, http.StatusTemporaryRedirect)
		}
		return
	}
	// error is ignored, because the id was checked in the AuthHandler
	sessionToken, err := r.Cookie(auth.SessionTokenKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := requests.GetRequestAuth[entities.Profile]("/v1/profile", sessionToken.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Profile")
		w.Header().Set("HX-Push-Url", "/profile")
		pages.Profile(user).Render(r.Context(), w)
		return
	}
	layouts.Default(layouts.PageProps{
		Title:       "Profile",
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/profile",
		Type:        layouts.ProfilePage,
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Profile(user)).Render(r.Context(), w)
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

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Results for "+query)
		w.Header().Set("HX-Push-Url", "/search?query="+query)
		pages.SearchResults(results).Render(r.Context(), w)
		return
	}
	layouts.Default(layouts.PageProps{
		Title:       "Results for " + query,
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/search?query=" + query,
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.SearchResults(results)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw(layouts.PageProps{
		Title:       "Signup",
		Description: "", // TODO:??
		Url:         "https://" + config.Env().Hostname + "/signup",
		ImageUrl:    "https://" + config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Signup()).Render(r.Context(), w)
}
