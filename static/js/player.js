"use strict";

// collapsed player's elements
const playPauseToggleEl = document.getElementById("play"),
  shuffleEl = document.getElementById("shuffle"),
  nextEl = document.getElementById("next"),
  prevEl = document.getElementById("prev"),
  loopEl = document.getElementById("loop"),
  songNameEl = document.getElementById("song-name"),
  artistNameEl = document.getElementById("artist-name"),
  songSeekBarEl = document.getElementById("song-seek-bar"),
  volumeSeekBarEl = document.getElementById("volume-seek-bar"),
  songDurationEl = document.getElementById("song-duration"),
  songCurrentTimeEl = document.getElementById("song-current-time"),
  songImageEl = document.getElementById("song-image"),
  audioPlayerEl = document.getElementById("audio-player"),
  muzikkContainerEl = document.getElementById("muzikk"),
  playerEl = document.getElementById("ze-player"),
  collapsedMobilePlayer = document.getElementById("ze-collapsed-mobile-player");

// expanded player's elements
const playPauseToggleExapndedEl = document.getElementById("play-expand"),
  songNameExpandedEl = document.getElementById("song-name-expanded"),
  artistNameExpandedEl = document.getElementById("artist-name-expanded"),
  songSeekBarExpandedEl = document.getElementById("song-seek-bar-expanded"),
  volumeSeekBarExpandedEl = document.getElementById("volume-seek-bar-expanded"),
  songDurationExpandedEl = document.getElementById("song-duration-expanded"),
  songCurrentTimeExpandedEl = document.getElementById(
    "song-current-time-expanded",
  ),
  songImageExpandedEl = document.getElementById("song-image-expanded"),
  expandedMobilePlayer = document.getElementById("ze-expanded-mobile-player");

/**
 * @typedef {object} Song
 * @property {string} title
 * @property {string} artist
 * @property {string} duration
 * @property {string} thumbnail_url
 * @property {string} yt_id
 * @property {number} play_times
 * @property {string} added_at
 */

/**
 * @typedef {object} Playlist
 * @property {string} public_id
 * @property {string} title
 * @property {string} songs_count
 * @property {Song[]} songs
 */

/**
 * @typedef {object} PlayerState
 * @property {LoopMode} loopMode
 * @property {boolean} shuffled
 * @property {Playlist} playlist
 * @property {number} currentSongIdx
 */

/**
 * @enum {LoopMode}
 */
const LOOP_MODES = Object.freeze({
  ALL: "ALL",
  OFF: "OFF",
  ONCE: "ONCE",
});

/**
 * @type{PlayerState}
 */
const playerState = {
  loopMode: LOOP_MODES.OFF,
  shuffled: false,
  currentSongIdx: 0,
  playlist: {
    title: "Queue",
    songs_count: 0,
    public_id: "",
    songs: [],
  },
};

function isSingleSong() {
  return playerState.playlist.songs.length <= 1;
}

/**
 * @returns {[Function, Function]}
 */
function looper() {
  const loopModes = [LOOP_MODES.OFF, LOOP_MODES.ONCE, LOOP_MODES.ALL];
  let currentLoopIdx = 0;

  const __toggle = () => {
    if (isSingleSong()) {
      currentLoopIdx = currentLoopIdx === 0 ? 1 : 0;
    } else {
      currentLoopIdx = (currentLoopIdx + 1) % loopModes.length;
    }
    loopEl.innerHTML =
      Player.icons[
        loopModes[currentLoopIdx] === LOOP_MODES.ALL
          ? "loop"
          : loopModes[currentLoopIdx] === LOOP_MODES.ONCE
            ? "loopOnce"
            : loopModes[currentLoopIdx] === LOOP_MODES.OFF
              ? "loopOff"
              : "loopOff"
      ];
  };

  const __handle = () => {
    switch (loopModes[currentLoopIdx]) {
      case LOOP_MODES.OFF:
        stopMuzikk();
        if (!isSingleSong()) {
          nextMuzikk();
        }
        break;
      case LOOP_MODES.ONCE:
        stopMuzikk();
        playMuzikk();
        break;
      case LOOP_MODES.ALL:
        if (!isSingleSong()) {
          nextMuzikk();
        }
        break;
    }
  };

  /**
   * @param {LoopMode} loopMode
   * @returns {boolean}
   */
  const __check = (loopMode) => loopMode === loopModes[currentLoopIdx];

  return [__toggle, __handle, __check];
}
/**
 * @param {HTMLElement} el
 * @param {string} icon
 */
const setPlayerButtonIcon = (el, icon) => {
  if (!!el && !!icon) {
    el.innerHTML = icon;
  }
};

/**
 * @param {boolean} loading
 * @param {string} fallback is used when loading is false, that is to reset
 *     the loading thingy
 */
function setLoading(loading, fallback) {
  if (loading) {
    setPlayerButtonIcon(playPauseToggleEl, Player.icons.loading);
    setPlayerButtonIcon(playPauseToggleExapndedEl, Player.icons.loading);
    document.body.style.cursor = "progress";
    return;
  }
  if (fallback) {
    setPlayerButtonIcon(playPauseToggleEl, fallback);
    setPlayerButtonIcon(playPauseToggleExapndedEl, fallback);
  }
  document.body.style.cursor = "auto";
}

/**
 * @param {HTMLAudioElement} audioEl
 *
 * @returns {[Function, Function, Function]}
 */
function playPauser(audioEl) {
  let startedPlaylist = false;

  const __play = () => {
    audioEl.play();
    const songEl = document.getElementById(
      "song-" + playerState.playlist.songs[playerState.currentSongIdx].yt_id,
    );
    if (!!songEl) {
      songEl.style.backgroundColor = "var(--accent-color-30)";
    }
    setPlayerButtonIcon(playPauseToggleEl, Player.icons.pause);
    setPlayerButtonIcon(playPauseToggleExapndedEl, Player.icons.pause);
  };
  const __pause = () => {
    audioEl.pause();
    setPlayerButtonIcon(playPauseToggleEl, Player.icons.play);
    setPlayerButtonIcon(playPauseToggleExapndedEl, Player.icons.play);
  };
  const __toggle = () => {
    const playPlaylistEl = document.getElementById("play-playlist-button");
    if (!!playPlaylistEl && !startedPlaylist) {
      playPlaylistEl.click();
      startedPlaylist = true;
      return;
    }
    if (audioEl.paused) {
      __play();
    } else {
      __pause();
    }
  };
  return [__play, __pause, __toggle];
}

/**
 * @param {HTMLAudioElement} audioEl
 *
 * @returns {Function}
 */
function stopper(audioEl) {
  return () => {
    audioEl.pause();
    audioEl.currentTime = 0;
    const songEl = document.getElementById(
      "song-" + playerState.playlist.songs[playerState.currentSongIdx].yt_id,
    );
    if (!!songEl) {
      songEl.style.backgroundColor = "#ffffff00";
    }
    setPlayerButtonIcon(playPauseToggleEl, Player.icons.play);
    setPlayerButtonIcon(playPauseToggleExapndedEl, Player.icons.play);
  };
}

/**
 * @param {PlayerState} state
 *
 * @returns {Function}
 */
function shuffler(state) {
  return () => {
    state.shuffled = !state.shuffled;
    setPlayerButtonIcon(
      shuffleEl,
      state.shuffled ? Player.icons.shuffle : Player.icons.shuffleOff,
    );
  };
}

/**
 * @param {PlayerState} state
 *
 * @returns {[Function, Function, Function, Function]}
 */
function playlister(state) {
  /**
   * @param {string} songYtId
   * @param {Playlist} playlist
   */
  const __setSongInPlaylistStyle = (songYtId, playlist) => {
    for (const _song of playlist.songs) {
      const songEl = document.getElementById("song-" + _song.yt_id);
      if (!songEl) {
        continue;
      }
      if (songYtId === _song.yt_id) {
        songEl.style.backgroundColor = "var(--accent-color-30)";
        songEl.scrollIntoView();
      } else {
        songEl.style.backgroundColor = "#ffffff00";
      }
    }
  };

  const __updateSongPlays = async () => {
    if (!state.playlist.public_id) {
      return;
    }
    await fetch(
      "/api/song/playlist/plays?" +
        new URLSearchParams({
          "song-id": state.playlist.songs[state.currentSongIdx].yt_id,
          "playlist-id": state.playlist.public_id,
        }).toString(),
      {
        method: "PUT",
      },
    ).catch((err) => console.error(err));
  };

  const __next = () => {
    if (checkLoop(LOOP_MODES.ONCE)) {
      stopMuzikk();
      playMuzikk();
      return;
    }
    if (
      !checkLoop(LOOP_MODES.ALL) &&
      !state.shuffled &&
      state.currentSongIdx + 1 >= state.playlist.songs.length
    ) {
      stopMuzikk();
      return;
    }
    state.currentSongIdx = state.shuffled
      ? Math.floor(Math.random() * state.playlist.songs.length)
      : checkLoop(LOOP_MODES.ALL) &&
          state.currentSongIdx + 1 >= state.playlist.songs.length
        ? 0
        : state.currentSongIdx + 1;
    const songToPlay = state.playlist.songs[state.currentSongIdx];
    playSongFromPlaylist(songToPlay.yt_id, state.playlist);
    __updateSongPlays();
    __setSongInPlaylistStyle(songToPlay.yt_id, state.playlist);
  };

  const __prev = () => {
    if (checkLoop(LOOP_MODES.ONCE)) {
      stopMuzikk();
      playMuzikk();
      return;
    }
    if (
      !checkLoop(LOOP_MODES.ALL) &&
      !state.shuffled &&
      state.currentSongIdx - 1 < 0
    ) {
      stopMuzikk();
      return;
    }
    state.currentSongIdx = state.shuffled
      ? Math.floor(Math.random() * state.playlist.songs.length)
      : checkLoop(LOOP_MODES.ALL) && state.currentSongIdx - 1 < 0
        ? state.playlist.songs.length - 1
        : state.currentSongIdx - 1;
    const songToPlay = state.playlist.songs[state.currentSongIdx];
    playSongFromPlaylist(songToPlay.yt_id, state.playlist);
    __updateSongPlays();
    __setSongInPlaylistStyle(songToPlay.yt_id, state.playlist);
  };

  const __remove = (songYtId, playlistId) => {
    const songIndex = state.playlist.songs.findIndex(
      (song) => song.yt_id === songYtId,
    );
    if (songIndex >= 0) {
      state.playlist.songs.splice(songIndex, 1);
    }

    Utils.showLoading();
    fetch(
      "/api/song/playlist?song-id=" +
        songYtId +
        "&playlist-id=" +
        playlistId +
        "&remove=true",
      {
        method: "PUT",
      },
    )
      .then((res) => {
        if (res.ok) {
          const songEl = document.getElementById("song-" + songYtId);
          if (!!songEl) {
            songEl.remove();
          }
        } else {
          window.alert("Oopsie something went wrong!");
        }
      })
      .catch((err) => {
        window.alert("Oopsie something went wrong!\n", err);
      })
      .finally(() => {
        Utils.hideLoading();
      });
  };

  return [__next, __prev, __remove, __setSongInPlaylistStyle];
}

function volumer() {
  let lastVolume = 1;
  const __setVolume = (level) => {
    if (level > 1) {
      level = 1;
    }
    if (level < 0) {
      level = 0;
    }
    audioPlayerEl.volume = level;
    if (volumeSeekBarEl) {
      volumeSeekBarEl.value = Math.floor(level * 100);
    }
    if (volumeSeekBarExpandedEl) {
      volumeSeekBarExpandedEl.value = Math.floor(level * 100);
    }
  };

  const __muter = () => {
    if (audioPlayerEl.volume === 0) {
      __setVolume(lastVolume);
    } else {
      lastVolume = audioPlayerEl.volume;
      __setVolume(0);
    }
  };

  return [__setVolume, __muter];
}

/**
 * @param {string} songYtId
 */
async function downloadSong(songYtId) {
  return await fetch("/api/song?id=" + songYtId).catch((err) =>
    console.error(err),
  );
}

/**
 * @param {string} songYtId
 * @param {string} songTitle
 */
async function downloadSongToDevice(songYtId, songTitle) {
  Utils.showLoading();
  await downloadSong(songYtId)
    .then(() => {
      const a = document.createElement("a");
      a.href = `/muzikkx/${songYtId}.mp3`;
      a.download = `${songTitle}.mp3`;
      a.click();
    })
    .finally(() => {
      Utils.hideLoading();
    });
}

/**
 * @param {string} songYtId
 */
async function downloadToApp() {
  throw new Error("not implemented!");
}

function show() {
  muzikkContainerEl.style.display = "block";
}

function hide() {
  muzikkContainerEl.style.display = "none";
}

function expand() {
  if (!playerEl.classList.contains("exapnded")) {
    playerEl.classList.add("exapnded");
    collapsedMobilePlayer.classList.add("hidden");
    expandedMobilePlayer.classList.remove("hidden");
  }
}

function collapse() {
  if (playerEl.classList.contains("exapnded")) {
    playerEl.classList.remove("exapnded");
    collapsedMobilePlayer.classList.remove("hidden");
    expandedMobilePlayer.classList.add("hidden");
  }
}

/**
 * @param {Song} song
 */
async function playSong(song) {
  setLoading(true);
  show();

  await downloadSong(song.yt_id).then(() => {
    stopMuzikk();
    audioPlayerEl.src = `/muzikkx/${song.yt_id}.mp3`;
    audioPlayerEl.load();
  });

  // song's details setting, yada yada
  {
    if (song.title) {
      songNameEl.innerHTML = song.title;
      songNameEl.title = song.title;
      if (song.title.length > Utils.getTextWidth()) {
        songNameEl.parentElement.classList.add("marquee");
      } else {
        songNameEl.parentElement.classList.remove("marquee");
      }

      if (songNameExpandedEl) {
        songNameExpandedEl.innerHTML = song.title;
        songNameExpandedEl.title = song.title;
        if (song.title.length > Utils.getTextWidth()) {
          songNameExpandedEl.parentElement.classList.add("marquee");
        } else {
          songNameExpandedEl.parentElement.classList.remove("marquee");
        }
      }
    }
    if (song.artist) {
      if (!!artistNameEl) {
        artistNameEl.innerHTML = song.artist;
        artistNameEl.title = song.artist;
      }

      if (artistNameExpandedEl) {
        artistNameExpandedEl.innerHTML = song.artist;
        artistNameExpandedEl.title = song.artist;
      }
    }
    songImageEl.style.backgroundImage = `url("${song.thumbnail_url}")`;
    songImageEl.innerHTML = "";

    if (songImageExpandedEl) {
      songImageExpandedEl.style.backgroundImage = `url("${song.thumbnail_url}")`;
      songImageExpandedEl.innerHTML = "";
    }
  }
  setMediaSessionMetadata(song);
  playMuzikk();
}

/**
 * @param {Song} song
 */
function playSingleSong(song) {
  playerState.playlist = {
    title: "Queue",
    songs_count: 1,
    public_id: "",
    songs: [song],
  };
  playerState.currentSongIdx = 0;
  playSong(song);
}

/**
 * @param {string} songYtId
 * @param {Playlist} playlist
 */
function playSongFromPlaylist(songYtId, playlist) {
  const songIdx = playlist.songs.findIndex((s) => s.yt_id === songYtId);
  if (songIdx < 0) {
    alert("Invalid song!");
    return;
  }
  playerState.playlist = playlist;
  playerState.currentSongIdx = songIdx;
  const songToPlay = playlist.songs[songIdx];
  highlightSongInPlaylist(songToPlay.yt_id, playlist);
  playSong(songToPlay);
}

/**
 * @param {Song} song
 */
function appendSongToCurrentQueue(song) {
  if (
    playerState.playlist.songs.findIndex((s) => s.yt_id === song.yt_id) !== -1
  ) {
    alert(`${song.title} exists in the queue!`);
    return;
  }
  playerState.playlist.songs.push(song);
  alert(`Added ${song.title} to the queue!`);
}

function addSongToPlaylist() {
  throw new Error("not implemented!");
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

const [toggleLoop, handleLoop, checkLoop] = looper();
const [playMuzikk, pauseMuzikk, togglePP] = playPauser(audioPlayerEl);
const stopMuzikk = stopper(audioPlayerEl);
const toggleShuffle = shuffler(playerState);
const [
  nextMuzikk,
  previousMuzikk,
  removeSongFromPlaylist,
  highlightSongInPlaylist,
] = playlister(playerState);
const [setVolume, mute] = volumer();

playPauseToggleEl.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  togglePP();
});

playPauseToggleExapndedEl?.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  togglePP();
});

nextEl?.addEventListener("click", nextMuzikk);
prevEl?.addEventListener("click", previousMuzikk);
shuffleEl?.addEventListener("click", toggleShuffle);

loopEl?.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  toggleLoop();
});

(() => {
  const __handler = (e) => {
    e.stopImmediatePropagation();
    e.preventDefault();
    const seekTime = Number(e.target.value);
    audioPlayerEl.currentTime = seekTime;
  };
  for (const event of ["change", "click"]) {
    songSeekBarEl?.addEventListener(event, __handler);
    songSeekBarExpandedEl?.addEventListener(event, __handler);
  }
})();

(() => {
  const __handler = (e) => {
    e.stopImmediatePropagation();
    e.preventDefault();
    const volume = Number(e.target.value) * 0.01;
    audioPlayerEl.volume = volume;
  };
  for (const event of ["change", "click"]) {
    volumeSeekBarEl?.addEventListener(event, __handler);
    volumeSeekBarExpandedEl?.addEventListener(event, __handler);
  }
})();

audioPlayerEl.addEventListener("loadeddata", (event) => {
  playPauseToggleEl.disabled = null;
  if (!!playPauseToggleExapndedEl) {
    playPauseToggleExapndedEl.disabled = null;
  }
  shuffleEl.disabled = null;
  nextEl.disabled = null;
  prevEl.disabled = null;
  loopEl.disabled = null;

  // set duration AAA
  {
    let duration = event.target.duration;
    if (isNaN(duration)) {
      duration = 0;
    }
    songSeekBarEl.max = Math.ceil(duration);
    songSeekBarEl.value = 0;
    if (!!songSeekBarExpandedEl) {
      songSeekBarExpandedEl.max = Math.ceil(duration);
      songSeekBarExpandedEl.value = 0;
    }
    if (!!songDurationEl) {
      songDurationEl.innerHTML = Utils.formatTime(duration);
    }
    if (!!songDurationExpandedEl) {
      songDurationExpandedEl.innerHTML = Utils.formatTime(duration);
    }
  }

  setLoading(false, Player.icons.pause);
});

audioPlayerEl.addEventListener("timeupdate", (event) => {
  const currentTime = Math.floor(event.target.currentTime);
  if (songCurrentTimeEl) {
    songCurrentTimeEl.innerHTML = Utils.formatTime(currentTime);
  }
  if (songCurrentTimeExpandedEl) {
    songCurrentTimeExpandedEl.innerHTML = Utils.formatTime(currentTime);
  }
  if (songSeekBarEl) {
    songSeekBarEl.value = Math.ceil(currentTime);
  }
  if (songSeekBarExpandedEl) {
    songSeekBarExpandedEl.value = Math.ceil(currentTime);
  }
});

audioPlayerEl.addEventListener("ended", () => {
  handleLoop();
});

audioPlayerEl.addEventListener("progress", () => {
  console.log("downloading...");
});

document
  .getElementById("collapse-player-button")
  ?.addEventListener("click", (event) => {
    event.stopImmediatePropagation();
    event.preventDefault();
    collapse();
  });

(() => {
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

  navigator.mediaSession.setActionHandler("previoustrack", () => {
    previousMuzikk();
  });

  navigator.mediaSession.setActionHandler("nexttrack", () => {
    nextMuzikk();
  });

  navigator.mediaSession.setActionHandler("stop", () => {
    stopMuzikk();
  });
})();

window.Player = {};
window.Player.downloadSongToDevice = downloadSongToDevice;
window.Player.showPlayer = show;
window.Player.hidePlayer = hide;
window.Player.playSingleSong = playSingleSong;
window.Player.playSongFromPlaylist = playSongFromPlaylist;
window.Player.removeSongFromPlaylist = removeSongFromPlaylist;
window.Player.addSongToQueue = appendSongToCurrentQueue;
window.Player.stopMuzikk = stopMuzikk;
window.Player.expand = () => expand();
window.Player.collapse = () => collapse();
