package pages

import (
	"context"
	"dankmuzikk-web/actions"
	"dankmuzikk-web/config"
	dankerrors "dankmuzikk-web/errors"
	"dankmuzikk-web/handlers/middlewares/auth"
	"dankmuzikk-web/handlers/middlewares/contenttype"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/status"
	"dankmuzikk-web/views/layouts"
	"dankmuzikk-web/views/pages"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/a-h/templ"
)

type pagesHandler struct {
	usecases *actions.Actions
}

func New(usecases *actions.Actions) *pagesHandler {
	return &pagesHandler{
		usecases: usecases,
	}
}

func (p *pagesHandler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	var recentPlays []actions.Song
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if ok {
		var err error
		recentPlays, err = p.usecases.GetHistory(sessionToken, 1)
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
		Url:         config.Env().Hostname,
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Index(recentPlays)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw(layouts.PageProps{
		Title:       "Login",
		Description: "", // TODO:??
		Url:         config.Env().Hostname + "/login",
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Login()).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSignupPage(w http.ResponseWriter, r *http.Request) {
	layouts.Raw(layouts.PageProps{
		Title:       "Signup",
		Description: "", // TODO:??
		Url:         config.Env().Hostname + "/signup",
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Signup()).Render(r.Context(), w)
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
		Url:         config.Env().Hostname + "/about",
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.About()).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePrivacyPage(w http.ResponseWriter, r *http.Request) {
	layouts.Default(layouts.PageProps{
		Title:       "Privacy",
		Description: "", // TODO:??
		Url:         config.Env().Hostname + "/privacy",
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Privacy()).Render(r.Context(), w)
}

func (p *pagesHandler) HandlePlaylistsPage(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			BugsBunnyError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	playlists, err := p.usecases.GetAllPlaylists(sessionToken)
	if err != nil {
		log.Errorln(err)
		playlists = make([]actions.Playlist, 0)
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
		Url:         config.Env().Hostname + "/playlists",
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Playlists(playlists)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSinglePlaylistPage(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		status.
			BugsBunnyError("I'm not sure what you're trying to do here :)").
			Render(r.Context(), w)
		return
	}

	playlistPubId := r.PathValue("playlist_id")
	if playlistPubId == "" {
		status.
			BugsBunnyError("You need to provide a playlist id!").
			Render(r.Context(), w)
		return
	}

	playlist, err := p.usecases.GetSinglePlaylist(sessionToken, playlistPubId)
	htmxReq := contenttype.IsNoLayoutPage(r)
	switch {
	case errors.Is(err, dankerrors.ErrUnauthorizedToSeePlaylist):
		log.Errorln(err)
		if htmxReq {
			status.
				BugsBunnyError("You can't see this playlist! <br/> (don't snoop around other people's stuff or else!)").
				Render(r.Context(), w)
			return
		} else {
			layouts.Default(layouts.PageProps{
				Title: "Error",
			},
				status.
					BugsBunnyError("You can't see this playlist! <br/> (don't snoop around other people's stuff or else!)")).
				Render(r.Context(), w)
			return
		}
	case err != nil:
		log.Errorln(err)
		if htmxReq {
			status.
				BugsBunnyError("You can't see this playlist! <br/> (it might be John Cena)").
				Render(r.Context(), w)
			return
		} else {
			layouts.Default(layouts.PageProps{
				Title: "Error",
			},
				status.
					BugsBunnyError("You can't see this playlist! <br/> (it might be John Cena)")).
				Render(r.Context(), w)
			return
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
		Url:         config.Env().Hostname + "/playlist/" + playlist.PublicId,
		Type:        layouts.PlaylistPage,
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Playlist(playlist)).Render(ctx, w)
}

func (p *pagesHandler) HandleSingleSongPage(w http.ResponseWriter, r *http.Request) {
	sessionToken, _ := r.Context().Value(auth.CtxSessionTokenKey).(string)

	songId := r.PathValue("song_id")
	if songId == "" {
		status.
			BugsBunnyError("You need to provide a song id!").
			Render(r.Context(), w)
		return
	}

	song, err := p.usecases.GetSongMetadata(sessionToken, songId)
	if err != nil {
		status.
			BugsBunnyError("Song doesn't exist!").
			Render(r.Context(), w)
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", song.Title)
		w.Header().Set("HX-Push-Url", "/song/"+song.PublicId)
		pages.Song(song).Render(r.Context(), w)
		return
	}
	layouts.Default(layouts.PageProps{
		Title:       song.Title,
		Description: "", // TODO:??
		Url:         config.Env().Hostname + "/song/" + song.PublicId,
		Type:        layouts.SongPage,
		ImageUrl:    song.ThumbnailUrl,
		Audio: layouts.AudioProps{
			Url:      fmt.Sprintf("%s/muzikkx/%s.mp3", config.Env().CdnAddress, song.PublicId),
			Duration: song.Duration,
			Musician: song.Artist,
		},
	}, pages.Song(song)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(auth.CtxSessionTokenKey).(string)
	if !ok {
		if contenttype.IsNoLayoutPage(r) {
			w.Header().Set("HX-Redirect", "/")
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		return
	}

	profile, err := p.usecases.GetProfile(sessionToken)
	if err != nil {
		if contenttype.IsNoLayoutPage(r) {
			w.Header().Set("HX-Redirect", "/")
		} else {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		return
	}

	if contenttype.IsNoLayoutPage(r) {
		w.Header().Set("HX-Title", "Profile")
		w.Header().Set("HX-Push-Url", "/profile")
		pages.Profile(profile).Render(r.Context(), w)
		return
	}
	layouts.Default(layouts.PageProps{
		Title:       "Profile",
		Description: "", // TODO:??
		Url:         config.Env().Hostname + "/profile",
		Type:        layouts.ProfilePage,
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.Profile(profile)).Render(r.Context(), w)
}

func (p *pagesHandler) HandleSearchResultsPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results, err := p.usecases.SearchYouTube(query)
	if err != nil {
		status.
			BugsBunnyError("Oopsie doopsie your query didn't result anything :)").
			Render(r.Context(), w)
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
		Url:         config.Env().Hostname + "/search?query=" + query,
		ImageUrl:    config.Env().Hostname + "/static/favicon-32x32.png",
	}, pages.SearchResults(results)).Render(r.Context(), w)
}
