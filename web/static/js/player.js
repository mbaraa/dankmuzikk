"use strict";

const audioPlayerEl = document.getElementById("muzikk-player");

/**
 * @type {Requests}
 */
let requests;

/**
 * @type {PlayerState}
 */
let freshPlayerState;

let removeSongFromPlaylist;
let setVolume, mute;
let setPlaybackSpeed;

async function init() {
  requests = new Requests();
  freshPlayerState = new PlayerState({});
  await freshPlayerState.refreshState().catch((e) => {
    console.log(e);
  });

  handleUIEvents();
  handlePlayerElementEvents();
  handleMediaSessionEvents();

  removeSongFromPlaylist = () => requests.removeSongFromPlaylist;
  [setVolume, mute] = volumer();
  setPlaybackSpeed = playebackSpeeder();
}

/**
 * @typedef {object} Song
 * @property {string} title
 * @property {string} artist
 * @property {string} duration
 * @property {string} thumbnail_url
 * @property {string} public_id
 * @property {number} play_times
 * @property {string} added_at
 * @property {number} votes
 * @property {number} order
 * @property {boolean} fully_downloaded
 * @property {string} media_url
 */

/**
 * @typedef {object} Playlist
 * @property {string} public_id
 * @property {string} title
 * @property {string} songs_count
 * @property {Song[]} songs
 */

class Requests {
  /**
   * @param {string} songPublicId
   * @param {string} playlistPublicId
   */
  async removeSongFromPlaylist(songPublicId, playlistPublicId) {
    Utils.showLoading();
    fetch(
      "/api/playlist/song?song-id=" +
        songPublicId +
        "&playlist-id=" +
        playlistPublicId +
        "&remove=true",
      {
        method: "PUT",
      },
    )
      .then((res) => {
        if (!res.ok) {
          alert("Oopsie something went wrong!");
        }
      })
      .catch((err) => {
        alert("Oopsie something went wrong!\n", err);
      })
      .finally(() => {
        Utils.hideLoading();
      });
  }

  /**
   * @param {string} songPublicId
   * @param {string} playlistPublicId
   */
  async downloadSong(songPublicId, playlistPublicId) {
    Utils.showLoading();

    try {
      let song = null;
      let songMetadataFetched = false;

      for (let i = 0; i < 30; i++) {
        const resp = await fetch(`/api/song/single?id=${songPublicId}`);
        if (!resp.ok) {
          throw new Error("Something went wrong when fetching song's metadata");
        }
        song = await resp.json();
        songMetadataFetched = true;

        if (song.fully_downloaded) {
          break;
        }
        await Utils.sleep(1000);
      }

      if (!songMetadataFetched) {
        throw new Error(
          "Failed to fetch song metadata after multiple attempts.",
        );
      }

      if (!song.fully_downloaded) {
        console.warn("Song not fully downloaded after repeated checks.");
      }

      const songDetailsResp = await fetch(
        "/api/song?id=" +
          songPublicId +
          (!!playlistPublicId ? `&playlist-id=${playlistPublicId}` : ""),
      );
      const songDetails = await songDetailsResp.json();

      if (songDetails.media_url) {
        const a = document.createElement("a");
        a.href = songDetails.media_url.replace("muzikkx", "muzikkx-raw");
        a.click();
        a.remove();
        return { ok: true, ...songDetails };
      } else {
        throw new Error("No media URL found for the song.");
      }
    } catch (err) {
      console.error("An error occurred during download:", err);
      return { ok: false, error: err.message };
    } finally {
      Utils.hideLoading();
    }
  }

  async downloadToApp() {
    throw new Error("not implemented!");
  }

  /**
   * @param {string} playlistPublicId
   */
  async downloadPlaylist(playlistPublicId) {
    Utils.showLoading();
    await fetch(`/api/playlist/zip?playlist-id=${playlistPublicId}`)
      .then(async (res) => {
        if (!res.ok) {
          throw new Error(await res.text());
        }
        return res.json();
      })
      .then((res) => {
        const a = document.createElement("a");
        a.href = res["playlist_download_url"];
        // a.download = `${plTitle}.zip`;
        a.click();
        a.remove();
      })
      .finally(() => {
        Utils.hideLoading();
      });
  }
}

class PlayerState {
  /**
   * @type {boolean}
   */
  #shuffled;
  /**
   * @type {number}
   */
  #currentSongIndex;
  /**
   * @type {"off" | "once" | "all"}
   */
  #loopMode;
  /**
   * @type {Song[]}
   */
  #songs;
  /**
   * @type {string}
   */
  #playingPlaylistPublicId;

  constructor(data) {
    this.#init(data);
  }

  #init(data) {
    this.#shuffled = data.shuffled ?? false;
    this.#currentSongIndex = data.current_song_index ?? 0;
    this.#loopMode = data.loop_mode ?? "off";
    this.#songs = data.songs ?? [];
    this.#playingPlaylistPublicId = data.playing_playlist_public_id ?? "";
  }

  async refreshState() {
    const result = await fetch("/api/player")
      .then((r) => r.json())
      .catch((e) => {
        console.error(e);
        return null;
      });
    if (!result) {
      return;
    }

    this.#init({
      ...result.player_state,
      playing_playlist_public_id: this.#playingPlaylistPublicId,
    });
  }

  isSingleSong() {
    return !!this.#songs && this.#songs.length === 1;
  }

  async toggleShuffle() {
    if (this.#shuffled) {
      await fetch("/api/player/shuffle", { method: "DELETE" });
      this.#shuffled = !this.#shuffled;
      PlayerUI.setShuffleOff();
    } else {
      await fetch("/api/player/shuffle", { method: "POST" });
      this.#shuffled = !this.#shuffled;
      PlayerUI.setShuffleOn();
    }
  }

  async toggleLoop() {
    const loopModes = ["off", "once", "all"];
    let loopModeIdx = loopModes.indexOf(this.#loopMode);
    if (this.isSingleSong()) {
      loopModeIdx = loopModeIdx === 0 ? 1 : 0;
    } else {
      loopModeIdx = (loopModeIdx + 1) % loopModes.length;
    }
    const loopMode = loopModes[loopModeIdx];
    switch (loopMode) {
      default:
      case "off":
        await fetch("/api/player/loop/off", { method: "PUT" });
        // TODO: global call aaaa
        PlayerUI.setLoopOff();
        this.#loopMode = "off";
        break;
      case "once":
        await fetch("/api/player/loop/once", { method: "PUT" });
        PlayerUI.setLoopOnce();
        this.#loopMode = "once";
        break;
      case "all":
        await fetch("/api/player/loop/all", { method: "PUT" });
        PlayerUI.setLoopAll();
        this.#loopMode = "all";
        break;
    }
  }

  play() {
    audioPlayerEl.muted = null;
    audioPlayerEl.play();
    PlayerUI.setPauseIcon();
  }

  pause() {
    audioPlayerEl.pause();
    PlayerUI.setPlayIcon();
  }

  togglePlayPause() {
    if (audioPlayerEl.paused) {
      this.play();
    } else {
      this.pause();
    }
  }

  stop() {
    audioPlayerEl.pause();
    audioPlayerEl.currentTime = 0;
    PlayerUI.setPlayIcon();
  }

  async next() {
    return await this.#nextPrevious(true);
  }

  async previous() {
    return await this.#nextPrevious(false);
  }

  async #nextPrevious(next = true) {
    /**
     * @type {{song: Song, current_song_index: number, end_of_queue: boolean}}
     */
    const result = await fetch(`/api/player/song/${next ? "next" : "previous"}`)
      .then((r) => r.json())
      .then((r) => r)
      .catch((e) => {
        console.error(e);
        return null;
      });
    if (!result) {
      return;
    }

    if (result.end_of_queue) {
      this.stop();
      return;
    }

    console.log("next song", result);

    // PlayerUI.highlightSongInPlaylist(
    // result.song.public_id,
    // freshPlayerState.songs.map((s) => s.public_id),
    // );

    return await this.fetchAndPlaySong(
      result.song.public_id,
      this.#playingPlaylistPublicId,
    );
  }

  /**
   * @param {string} songPublicId
   */
  async #fetchSongMetadata(songPublicId, displayLoader = true) {
    if (displayLoader) {
      Utils.showLoading();
    }
    try {
      const song = await fetch(`/api/song/single?id=${songPublicId}`);
      if (!song.ok) {
        throw new Error("Something went wrong when fetching song's metadata");
      }
      return await song.json();
    } catch (err) {
      return err;
    } finally {
      if (displayLoader) {
        Utils.hideLoading();
      }
    }
  }

  /**
   * @param {string} songPublicId
   * @param {string} playlistPublicId
   */
  async fetchAndPlaySong(songPublicId, playlistPublicId) {
    let newState = false;
    if (playlistPublicId === "") {
      this.#playingPlaylistPublicId = "";
      newState = true;
    }
    if (playlistPublicId !== this.#playingPlaylistPublicId) {
      this.#playingPlaylistPublicId = playlistPublicId;
      newState = true;
    }

    console.log("new state", newState);

    PlayerUI.setLoadingOn();

    const resp = await fetch(
      "/api/song?id=" +
        songPublicId +
        (!!playlistPublicId ? `&playlist-id=${playlistPublicId}` : ""),
    )
      .then((res) => res.json())
      .then(async (data) => {
        for (let i = 0; i < 15; i++) {
          const song = await this.#fetchSongMetadata(songPublicId, false);
          if (song.fully_downloaded) {
            return { ...data, ...song };
          }
          await Utils.sleep(1000);
        }
      })
      .catch((e) => {
        console.error(e);
        return null;
      });
    if (!resp) {
      alert("Something went wrong when downloading the song...");
    }

    this.stop();
    if (audioPlayerEl.childNodes.length > 0) {
      audioPlayerEl.removeChild(audioPlayerEl.childNodes.item(0));
    }
    const src = document.createElement("source");
    src.setAttribute("type", "audio/mpeg");
    src.setAttribute("src", resp.media_url);
    src.setAttribute("preload", "metadata");
    audioPlayerEl.appendChild(src);
    audioPlayerEl.load();

    PlayerUI.setSongName(resp.title);
    PlayerUI.setArtistName(resp.artist);
    PlayerUI.setSongThumbnail(resp.thumbnail_url);
    setMediaSessionMetadata(resp);
    this.play();

    if (newState) {
      await this.refreshState();
    }
  }

  async loadSongLyrics() {
    document.getElementById("current-song-lyrics").innerHTML = "";
    const songPublicId = this.#songs[this.#currentSongIndex].public_id;
    htmx
      .ajax("GET", "/api/song/lyrics?id=" + songPublicId, {
        target: "#current-song-lyrics",
        swap: "innerHTML",
      })
      .catch(() => {
        alert("Lyrics fetching went berzerk...");
      });
  }

  /**
   * @param {string} songPublicId
   */
  async addSongToQueueNext(songPublicId) {
    await fetch(`/api/player/queue/song/next?id=${songPublicId}`, {
      method: "POST",
    }).catch((e) => {
      console.error(e);
    });

    await this.refreshState();
  }

  /**
   * @param {string} songPublicId
   */
  async addSongToQueueAtLast(songPublicId) {
    await fetch(`/api/player/queue/song/last?id=${songPublicId}`, {
      method: "POST",
    }).catch((e) => {
      console.error(e);
    });

    await this.refreshState();
  }

  /**
   * @param {string} playlistPublicId
   */
  async addPlaylistToQueueNext(playlistPublicId) {
    await fetch(`/api/player/queue/playlist/next?id=${playlistPublicId}`, {
      method: "POST",
    }).catch((e) => {
      console.error(e);
    });

    await this.refreshState();
  }

  /**
   * @param {string} playlistPublicId
   */
  async addPlaylistToQueueAtLast(playlistPublicId) {
    await fetch(`/api/player/queue/playlist/last?id=${playlistPublicId}`, {
      method: "POST",
    }).catch((e) => {
      console.error(e);
    });

    await this.refreshState();
  }
}

function volumer() {
  const __setVolume = (level) => {
    if (level > 1) {
      level = 1;
    }
    if (level < 0) {
      level = 0;
    }
    audioPlayerEl.volume = level;
    PlayerUI.setVolumeLevel(level);
  };

  const __muter = () => {
    audioPlayerEl.muted = !audioPlayerEl.muted;
  };

  return [__setVolume, __muter];
}

function playebackSpeeder() {
  /**
   * @param {number} speed
   */
  const __setSpeed = (speed) => {
    speed = speed < 0.1 ? 0.1 : speed > 4 ? 4 : speed;
    audioPlayerEl.playbackRate = speed;
    // TODO: add the ui stuff
  };

  return __setSpeed;
}

/**
 * @param {string} songPublicId
 * @param {string} songTitle
 */
async function downloadSongToDevice(songPublicId, songTitle) {
  return await requests.downloadSong(songPublicId);
}

/**
 * @param {string} plPubId
 * @param {plTitle} plTitle
 */
async function downloadPlaylistToDevice(plPubId, plTitle) {
  return await requests.downloadPlaylist(plPubId);
}

/**
 * @param {Song} song
 */
async function playSingleSong(song) {
  await window.Utils.retryer(async () => {
    return await freshPlayerState.fetchAndPlaySong(song.public_id);
  });
}

/**
 * @param {string} songPublicId
 */
async function playSingleSongId(songPublicId) {
  await window.Utils.retryer(async () => {
    return await freshPlayerState.fetchAndPlaySong(songPublicId);
  });
}

/**
 * @param {Song} song
 */
async function playSingleSongNext(song) {
  return await freshPlayerState.addSongToQueueNext(song.public_id);
}

/**
 * @param {string} songPublicId
 */
async function playSingleSongNextId(songPublicId) {
  return await freshPlayerState.addSongToQueueNext(songPublicId);
}

/**
 * @param {Playlist} playlist
 */
async function playPlaylistNext(playlist) {
  return await freshPlayerState.addPlaylistToQueueNext(playlist.public_id);
}

/**
 * @param {string} playlistPubId
 */
async function playPlaylistNextId(playlistPubId) {
  return await freshPlayerState.addPlaylistToQueueNext(playlistPubId);
}

/**
 * @param {Playlist} playlist
 */
async function appendPlaylistToCurrentQueue(playlist) {
  return await freshPlayerState.addPlaylistToQueueAtLast(playlist.public_id);
}

/**
 * @param {string} playlistPubId
 */
async function appendPlaylistToCurrentQueueId(playlistPubId) {
  return await freshPlayerState.addPlaylistToQueueAtLast(playlistPubId);
}

/**
 * @param {string} songPublicId
 * @param {Playlist} playlist
 */
async function playSongFromPlaylist(songPublicId, playlist) {
  return await freshPlayerState.fetchAndPlaySong(
    songPublicId,
    playlist.public_id,
  );
}

/**
 * @param {string} songPublicId
 * @param {string} playlistPubId
 */
async function playSongFromPlaylistId(songPublicId, playlistPubId) {
  return await freshPlayerState.fetchAndPlaySong(songPublicId, playlistPubId);
}

/**
 * @param {Song} song
 */
async function appendSongToCurrentQueue(song) {
  return await freshPlayerState.addSongToQueueAtLast(song.public_id);
}

/**
 * @param {string} songPublicId
 */
async function appendSongToCurrentQueueId(songPublicId) {
  return await freshPlayerState.addSongToQueueAtLast(songPublicId);
}

function playMuzikk() {
  return freshPlayerState.play();
}

function pauseMuzikk() {
  return freshPlayerState.pause();
}

function togglePlayPause() {
  return freshPlayerState.togglePlayPause();
}

function stopMuzikk() {
  return freshPlayerState.stop();
}

async function nextMuzikk() {
  return await freshPlayerState.next();
}

async function previousMuzikk() {
  return await freshPlayerState.previous();
}

async function toggleLoop() {
  return await freshPlayerState.toggleLoop();
}

async function toggleShuffle() {
  return await freshPlayerState.toggleShuffle();
}

async function loadSongLyrics() {
  await freshPlayerState.loadSongLyrics();
}

function handleUIEvents() {
  PlayerUI.__elements.playPauseToggleEl.addEventListener("click", (event) => {
    event.stopImmediatePropagation();
    event.preventDefault();
    togglePlayPause();
  });

  PlayerUI.__elements.playPauseToggleExapndedEl?.addEventListener(
    "click",
    (event) => {
      event.stopImmediatePropagation();
      event.preventDefault();
      togglePlayPause();
    },
  );

  PlayerUI.__elements.nextEl?.addEventListener("click", nextMuzikk);
  PlayerUI.__elements.nextExpandEl?.addEventListener("click", nextMuzikk);
  PlayerUI.__elements.prevEl?.addEventListener("click", previousMuzikk);
  PlayerUI.__elements.prevExpandEl?.addEventListener("click", previousMuzikk);
  PlayerUI.__elements.shuffleEl?.addEventListener("click", toggleShuffle);
  PlayerUI.__elements.shuffleExpandEl?.addEventListener("click", toggleShuffle);
  PlayerUI.__elements.loopEl?.addEventListener("click", toggleLoop);
  PlayerUI.__elements.loopExpandEl?.addEventListener("click", toggleLoop);

  {
    const __handler = (e) => {
      e.stopImmediatePropagation();
      e.preventDefault();
      const seekTime = Number(e.target.value);
      audioPlayerEl.currentTime = seekTime;
    };
    for (const event of ["change", "click"]) {
      PlayerUI.__elements.songSeekBarEl?.addEventListener(event, __handler);
      PlayerUI.__elements.songSeekBarExpandedEl?.addEventListener(
        event,
        __handler,
      );
    }
  }

  {
    const __handler = (e) => {
      e.stopImmediatePropagation();
      e.preventDefault();
      const volume = Number(e.target.value) * 0.01;
      audioPlayerEl.volume = volume;
    };
    for (const event of ["change", "click"]) {
      PlayerUI.__elements.volumeSeekBarEl?.addEventListener(event, __handler);
      PlayerUI.__elements.volumeSeekBarExpandedEl?.addEventListener(
        event,
        __handler,
      );
    }
  }
}

function handlePlayerElementEvents() {
  audioPlayerEl.addEventListener("loadeddata", (event) => {
    PlayerUI.enableButtons();
    PlayerUI.setSongDuration(event.target.duration);

    PlayerUI.setLoadingOff();
  });

  audioPlayerEl.addEventListener("timeupdate", (event) => {
    PlayerUI.setSongCurrentTime(event.target.currentTime);
  });

  audioPlayerEl.addEventListener("ended", () => {
    nextMuzikk();
  });

  audioPlayerEl.addEventListener("progress", () => {
    console.log("downloading...");
  });
}

function handleMediaSessionEvents() {
  if (!("mediaSession" in navigator)) {
    console.error("Browser doesn't support mediaSession");
    return;
  }

  navigator.mediaSession.setActionHandler("play", () => {
    playMuzikk();
  });

  navigator.mediaSession.setActionHandler("pause", () => {
    pauseMuzikk();
  });

  navigator.mediaSession.setActionHandler("stop", () => {
    stopMuzikk();
  });

  navigator.mediaSession.setActionHandler("seekbackward", () => {
    let seekTo = -10;
    if (audioPlayerEl.currentTime + seekTo < 0) {
      seekTo = 0;
    }
    audioPlayerEl.currentTime += seekTo;
  });

  navigator.mediaSession.setActionHandler("seekforward", () => {
    let seekTo = +10;
    if (audioPlayerEl.currentTime + seekTo > audioPlayerEl.duration) {
      seekTo = 0;
    }
    audioPlayerEl.currentTime += seekTo;
  });

  navigator.mediaSession.setActionHandler("seekto", (a) => {
    const seekTime = Number(a.seekTime);
    audioPlayerEl.currentTime = seekTime;
  });

  navigator.mediaSession.setActionHandler("previoustrack", async () => {
    await previousMuzikk();
  });

  navigator.mediaSession.setActionHandler("nexttrack", async () => {
    await nextMuzikk();
  });
}

/**
 * @param {Song} song
 */
function setMediaSessionMetadata(song) {
  if (!("mediaSession" in navigator)) {
    console.error("Browser doesn't support mediaSession");
    return;
  }

  navigator.mediaSession.metadata = new MediaMetadata({
    title: song.title,
    artist: song.artist,
    album: song.artist,
    artwork: [
      "96x96",
      "128x128",
      "192x192",
      "256x256",
      "384x384",
      "512x512",
    ].map((i) => {
      return {
        src: song.thumbnail_url,
        sizes: i,
        type: "image/png",
      };
    }),
  });
}

init();

window.Player = {
  downloadSongToDevice,
  downloadPlaylistToDevice,
  playSingleSong,
  playSingleSongId,
  playSingleSongNext,
  playSingleSongNextId,
  playSongFromPlaylist,
  playSongFromPlaylistId,
  playPlaylistNext,
  playPlaylistNextId,
  removeSongFromPlaylist,
  addSongToQueue: appendSongToCurrentQueue,
  addSongToQueueId: appendSongToCurrentQueueId,
  appendPlaylistToCurrentQueue,
  appendPlaylistToCurrentQueueId,
  setPlaybackSpeed,
  loadSongLyrics,
};
