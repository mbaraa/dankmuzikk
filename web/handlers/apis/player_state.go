package apis

import (
	"dankmuzikk-web/actions"
	"dankmuzikk-web/log"
	"dankmuzikk-web/views/components/lyrics"
	"dankmuzikk-web/views/components/player"
	"dankmuzikk-web/views/components/playlist"
	"dankmuzikk-web/views/components/song"
	"dankmuzikk-web/views/components/status"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type playerStateApi struct {
	usecases *actions.Actions
}

func NewPlayerStateApi(usecases *actions.Actions) *playerStateApi {
	return &playerStateApi{
		usecases: usecases,
	}
}

func (p *playerStateApi) HandleGetPlayerState(w http.ResponseWriter, r *http.Request) {
	ctx, _ := parseContext(r.Context())

	payload, err := p.usecases.GetPlayerState(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleGetPlayerSongsQueue(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := p.usecases.GetPlayerState(ctx)
	if err != nil {
		status.BugsBunnyError("No songs were found!\nMaybe play something first...").
			Render(r.Context(), w)
		return
	}

	for idx, s := range payload.PlayerState.Songs {
		song.Song(s, []string{s.AddedAt},
			[]templ.Component{
				song.RemoveFromQueue(s, idx),
				playlist.PlaylistsPopup((idx + 1), s.PublicId),
			},
			actions.Playlist{}, "queue").
			Render(r.Context(), w)
	}
}

func (p *playerStateApi) HandleSetPlayerShuffleOn(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	err = p.usecases.SetPlayerShuffleOn(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.ShuffleButton(true).Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerShuffleOff(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	err = p.usecases.SetPlayerShuffleOff(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.ShuffleButton(false).Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerLoopOff(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	err = p.usecases.SetPlayerLoopOff(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.LoopButton("off").Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerLoopOnce(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	err = p.usecases.SetPlayerLoopOnce(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.LoopButton("once").Render(r.Context(), w)
}

func (p *playerStateApi) HandleSetPlayerLoopAll(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	err = p.usecases.SetPlayerLoopAll(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	player.LoopButton("all").Render(r.Context(), w)
}

func (p *playerStateApi) HandleGetNextSongInQueue(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := p.usecases.GetNextSongInQueue(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleGetPreviousSongInQueue(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	payload, err := p.usecases.GetPreviousSongInQueue(ctx)
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(payload)
}

func (p *playerStateApi) HandleGetPlayingSongLyrics(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("You need to login to view cureent song's lyrics!").
			Render(r.Context(), w)
		return
	}

	lyricsResp, err := p.usecases.GetPlayingSongLyrics(ctx)
	if err != nil || len(lyricsResp.Lyrics) == 0 {
		status.BugsBunnyError("No Lyrics was found!").
			Render(r.Context(), w)
		return
	}

	_ = lyrics.Lyrics(lyricsResp.SongTitle, lyricsResp.Lyrics, lyricsResp.SyncedPairs()).
		Render(r.Context(), w)
}

func (p *playerStateApi) HandleAddSongToQueueNext(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	songId := r.URL.Query().Get("id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.usecases.AddSongToQueueNext(actions.AddSongToQueueNextParams{
		ActionContext: ctx,
		SongPublicId:  songId,
	})
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleAddSongToQueueAtLast(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	songId := r.URL.Query().Get("id")
	if songId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.usecases.AddSongToQueueAtLast(actions.AddSongToQueueNextParams{
		ActionContext: ctx,
		SongPublicId:  songId,
	})
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleRemoveSongFromQueue(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	songIndex, err := strconv.Atoi(r.URL.Query().Get("index"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.usecases.RemoveSongFromQueue(actions.RemoveSongFromQueueParams{
		ActionContext: ctx,
		SongIndex:     songIndex,
	})
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleAddPlaylistToQueueNext(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	playlistId := r.URL.Query().Get("id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.usecases.AddPlaylistToQueueNext(actions.AddPlaylistToQueueNextParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *playerStateApi) HandleAddPlaylistToQueueAtLast(w http.ResponseWriter, r *http.Request) {
	ctx, err := parseContext(r.Context())
	if err != nil {
		status.BugsBunnyError("What do you think you're doing?").
			Render(r.Context(), w)
		return
	}

	playlistId := r.URL.Query().Get("id")
	if playlistId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.usecases.AddPlaylistToQueueAtLast(actions.AddPlaylistToQueueAtLastParams{
		ActionContext:    ctx,
		PlaylistPublicId: playlistId,
	})
	if err != nil {
		log.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
