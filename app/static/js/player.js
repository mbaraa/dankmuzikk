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
  shuffleExpandEl = document.getElementById("shuffle-expand"),
  nextExpandEl = document.getElementById("next-expand"),
  prevExpandEl = document.getElementById("prev-expand"),
  loopExpandEl = document.getElementById("loop-expand"),
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
 * @property {number} votes
 * @property {number} order
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
 * @property {string} shuffledPlaylist
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
  shuffledPlaylist: "",
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
    const icon =
      Player.icons[
        loopModes[currentLoopIdx] === LOOP_MODES.ALL
          ? "loop"
          : loopModes[currentLoopIdx] === LOOP_MODES.ONCE
            ? "loopOnce"
            : loopModes[currentLoopIdx] === LOOP_MODES.OFF
              ? "loopOff"
              : "loopOff"
      ];
    setPlayerButtonIcon(loopEl, icon);
    setPlayerButtonIcon(loopExpandEl, icon);
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
  const __play = () => {
    audioEl.muted = null;
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
  // using Fisher–Yates shuffling algorithm https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
  const __shuffleArray = (a) => {
    let currIdx = a.length;
    while (currIdx != 0) {
      let randIdx = Math.floor(Math.random() * currIdx);
      currIdx--;
      [a[currIdx], a[randIdx]] = [a[randIdx], a[currIdx]];
    }
  };

  /**
   * @param {string} songYtId
   */
  const __shuffle = (songYtId) => {
    if (isSingleSong()) {
      return;
    }
    state.shuffledPlaylist = state.playlist.public_id;
    const extraSongs = [];
    state.playlist.songs.forEach((s) => {
      for (let i = 0; i < s.votes - 1; i++) {
        extraSongs.push(s);
      }
    });
    state.playlist.songs.push(...extraSongs);
    __shuffleArray(state.playlist.songs);
    let sIdx = 0;
    if (!!audioPlayerEl.childNodes.length) {
      sIdx = state.playlist.songs.findIndex((s) => s.yt_id === songYtId);
      if (sIdx !== -1) {
        [state.playlist.songs[sIdx], state.playlist.songs[0]] = [
          state.playlist.songs[0],
          state.playlist.songs[sIdx],
        ];
      }
    }
    state.currentSongIdx = 0;
  };

  const __toggleShuffle = () => {
    state.shuffled = !state.shuffled;
    if (state.shuffled) {
      const src = audioPlayerEl.childNodes.item(0);
      shuffle(
        src.src.substring(src.src.lastIndexOf("/") + 1, src.src.length - 4),
      );
    } else {
      const tmp = state.playlist.songs.sort((si, sj) => si.order - sj.order);
      state.playlist.songs = [];
      for (let i = 0; i < tmp.length - 1; i++) {
        if (tmp[i].yt_id === tmp[i + 1].yt_id) {
          continue;
        }
        state.playlist.songs.push(tmp[i]);
      }
      if (tmp[tmp.length - 1].yt_id !== tmp[tmp.length - 2]) {
        state.playlist.songs.push(tmp[tmp.length - 1]);
      }
      if (!!src) {
        console.log("src", src.src);
        state.currentSongIdx = state.playlist.songs.findIndex(
          (s) =>
            s.yt_id ===
            src.src.substring(src.src.lastIndexOf("/") + 1, src.src.length - 4),
        );
      }
    }
    setPlayerButtonIcon(
      shuffleEl,
      state.shuffled ? Player.icons.shuffle : Player.icons.shuffleOff,
    );
    setPlayerButtonIcon(
      shuffleExpandEl,
      state.shuffled ? Player.icons.shuffle : Player.icons.shuffleOff,
    );
  };

  return [__shuffle, __toggleShuffle];
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

  const __next = async () => {
    if (checkLoop(LOOP_MODES.ONCE)) {
      stopMuzikk();
      playMuzikk();
      await updateSongPlays();
      return;
    }
    // chack votes to whether repeat the song or not.
    if (
      state.playlist.songs[state.currentSongIdx].votes > 1 &&
      !state.shuffled
    ) {
      const songToPlay = state.playlist.songs[state.currentSongIdx];
      songToPlay.votes--;
      songToPlay.plays++;
      playSongFromPlaylist(songToPlay.yt_id, state.playlist);
      __setSongInPlaylistStyle(songToPlay.yt_id, state.playlist);
      return;
    }

    if (
      !checkLoop(LOOP_MODES.ALL) &&
      state.currentSongIdx + 1 >= state.playlist.songs.length
    ) {
      stopMuzikk();
      // reset songs' votes
      for (const s of state.playlist.songs) {
        if (!!s.plays) {
          s.votes = s.plays;
          s.plays = 0;
        }
      }
      return;
    }

    state.currentSongIdx =
      checkLoop(LOOP_MODES.ALL) &&
      state.currentSongIdx + 1 >= state.playlist.songs.length
        ? 0
        : state.currentSongIdx + 1;
    const songToPlay = state.playlist.songs[state.currentSongIdx];
    await highlightSongInPlaylist(songToPlay.yt_id, state.playlist);
    await playSong(songToPlay);
    __setSongInPlaylistStyle(songToPlay.yt_id, state.playlist);
  };

  const __prev = async () => {
    if (checkLoop(LOOP_MODES.ONCE)) {
      stopMuzikk();
      playMuzikk();
      await updateSongPlays();
      return;
    }
    // chack votes to whether repeat the song or not.
    if (
      state.playlist.songs[state.currentSongIdx].votes > 1 &&
      !state.shuffled
    ) {
      const songToPlay = state.playlist.songs[state.currentSongIdx];
      songToPlay.votes--;
      songToPlay.plays++;
      playSongFromPlaylist(songToPlay.yt_id, state.playlist);
      __setSongInPlaylistStyle(songToPlay.yt_id, state.playlist);
      return;
    }
    if (!checkLoop(LOOP_MODES.ALL) && state.currentSongIdx - 1 < 0) {
      stopMuzikk();
      // reset songs' votes
      for (const s of state.playlist.songs) {
        if (!!s.plays) {
          s.votes = s.plays;
          s.plays = 0;
        }
      }
      return;
    }
    state.currentSongIdx =
      checkLoop(LOOP_MODES.ALL) && state.currentSongIdx - 1 < 0
        ? state.playlist.songs.length - 1
        : state.currentSongIdx - 1;
    const songToPlay = state.playlist.songs[state.currentSongIdx];
    await highlightSongInPlaylist(songToPlay.yt_id, state.playlist);
    await playSong(songToPlay);
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

  return [__setSpeed];
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

/**
 * @param {string} plPubId
 * @param {plTitle} plTitle
 */
async function downloadPlaylistToDevice(plPubId, plTitle) {
  Utils.showLoading();
  await fetch(`/api/playlist/zip?playlist-id=${plPubId}`)
    .then(async (res) => {
      if (!res.ok) {
        throw new Error(await res.text());
      }
      return res.blob();
    })
    .then((playlistZip) => {
      const a = document.createElement("a");
      a.href = URL.createObjectURL(playlistZip);
      a.download = `${plTitle}.zip`;
      a.click();
    })
    .finally(() => {
      Utils.hideLoading();
    });
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

async function updateSongPlays() {
  if (!playerState.playlist.public_id) {
    return;
  }
  await fetch(
    "/api/song/playlist/plays?" +
      new URLSearchParams({
        "song-id": playerState.playlist.songs[playerState.currentSongIdx].yt_id,
        "playlist-id": playerState.playlist.public_id,
      }).toString(),
    {
      method: "PUT",
    },
  ).catch((err) => console.error(err));
}

/**
 * @param {Song} song
 */
async function playSong(song) {
  setLoading(true);
  show();

  const resp = await downloadSong(song.yt_id);
  if (!resp.ok) {
    alert("Something went wrong when downloading the song...");
    return;
  }
  stopMuzikk();
  if (audioPlayerEl.childNodes.length > 0) {
    audioPlayerEl.removeChild(audioPlayerEl.childNodes.item(0));
  }
  const src = document.createElement("source");
  src.src = `${location.protocol}//${location.host}/muzikkx/${song.yt_id}.mp3`;
  src.type = "audio/mpeg";
  audioPlayerEl.appendChild(src);

  if (isSafari()) {
    setTimeout(80);
  }
  audioPlayerEl.load();

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
  await updateSongPlays();
}

/**
 * @param {string} songYtId
 */
async function fetchSongMeta(songYtId) {
  Utils.showLoading();
  return await fetch(`/api/song/single?id=${songYtId}`)
    .then((res) => res.json())
    .then((s) => s)
    .catch((err) => {
      console.error(err);
    })
    .finally(() => {
      Utils.hideLoading();
    });
}

/**
 * @param {string} playlistPubId
 */
async function fetchPlaylistMeta(playlistPubId) {
  Utils.showLoading();
  return await fetch(`/api/playlist?playlist-id=${playlistPubId}`)
    .then((res) => res.json())
    .then((p) => p)
    .catch((err) => {
      console.error(err);
    })
    .finally(() => {
      Utils.hideLoading();
    });
}

/**
 * @param {Song} song
 */
async function playSingleSong(song) {
  playerState.playlist = {
    title: "Queue",
    songs_count: 1,
    public_id: "",
    songs: [song],
  };
  playerState.currentSongIdx = 0;
  await playSong(song);
}

/**
 * @param {string} songYtId
 */
async function playSingleSongId(songYtId) {
  const song = await fetchSongMeta(songYtId);
  await playSingleSong(song);
}

/**
 * @param {Song} song
 */
async function playSingleSongNext(song) {
  if (playerState.playlist.songs.length === 0) {
    playSingleSong(song);
    return;
  }
  if (!song.yt_id) {
    return;
  }
  song.votes = 1;
  playerState.playlist.songs.splice(playerState.currentSongIdx + 1, 0, song);
  alert(`Playing ${song.title} next!`);
}

/**
 * @param {string} songYtId
 */
async function playSingleSongNextId(songYtId) {
  const song = await fetchSongMeta(songYtId);
  await playSingleSongNext(song);
}

/**
 * @param {Playlist} playlist
 */
async function playPlaylistNext(playlist) {
  if (!playlist || !playlist.songs || playlist.songs.length === 0) {
    alert("Can't do that!");
    return;
  }
  if (playerState.playlist.songs.length === 0) {
    playSongFromPlaylist(playlist.songs[0].yt_id, playlist);
    return;
  }
  playerState.playlist.songs.splice(
    playerState.currentSongIdx + 1,
    0,
    ...playlist.songs.map((s) => {
      return { ...s, votes: 1 };
    }),
  );
  playerState.playlist.title = `${playerState.playlist.title} + ${playlist.title}`;
  alert(`Playing ${playlist.title} next!`);
}

/**
 * @param {string} playlistPubId
 */
async function playPlaylistNextId(playlistPubId) {
  const playlist = await fetchPlaylistMeta(playlistPubId);
  await playPlaylistNext(playlist);
}

/**
 * @param {Playlist} playlist
 */
async function appendPlaylistToCurrentQueue(playlist) {
  if (!playlist || !playlist.songs || playlist.songs.length === 0) {
    alert("Can't do that!");
    return;
  }
  if (playerState.playlist.songs.length === 0) {
    playSongFromPlaylist(playlist.songs[0].yt_id, playlist);
    return;
  }
  playerState.playlist.songs.push(...playlist.songs);
  playerState.playlist.title = `${playerState.playlist.title} + ${playlist.title}`;
  alert(`Playing ${playlist.title} next!`);
}

/**
 * @param {string} playlistPubId
 */
async function appendPlaylistToCurrentQueueId(playlistPubId) {
  const playlist = await fetchPlaylistMeta(playlistPubId);
  await appendPlaylistToCurrentQueue(playlist);
}

/**
 * @param {string} songYtId
 * @param {Playlist} playlist
 */
async function playSongFromPlaylist(songYtId, playlist) {
  if (
    playerState.shuffled &&
    playerState.shuffledPlaylist !== playlist.public_id
  ) {
    playerState.playlist = playlist;
    shuffle(songYtId);
  }
  const songIdx = playlist.songs.findIndex((s) => s.yt_id === songYtId);
  if (songIdx < 0) {
    alert("Invalid song!");
    return;
  }
  if (!playerState.shuffled) {
    playerState.playlist = playlist;
    playerState.playlist.songs = playlist.songs.map((s, idx) => {
      return { ...s, plays: 0, order: idx };
    });
  }
  playerState.currentSongIdx = songIdx;
  if (
    playerState.playlist.songs[songIdx].votes === 0 &&
    playerState.playlist.public_id !== ""
  ) {
    playerState.currentSongIdx++;
  }
  const songToPlay = playlist.songs[playerState.currentSongIdx];
  await highlightSongInPlaylist(songToPlay.yt_id, playlist);
  await playSong(songToPlay);
}

/**
 * @param {string} songYtId
 * @param {string} playlistPubId
 */
async function playSongFromPlaylistId(songYtId, playlistPubId) {
  const playlist = await fetchPlaylistMeta(playlistPubId);
  await playSongFromPlaylist(songYtId, playlist);
}

/**
 * @param {Song} song
 */
function appendSongToCurrentQueue(song) {
  if (playerState.playlist.songs.length === 0) {
    playSingleSong(song);
    return;
  }
  song.votes = 1;
  playerState.playlist.songs.push(song);
  alert(`Added ${song.title} to the queue!`);
}

/**
 * @param {string} songYtId
 */
async function appendSongToCurrentQueueId(songYtId) {
  const song = await fetchSongMeta(songYtId);
  appendSongToCurrentQueue(song);
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

function isSafari() {
  return navigator.userAgent.toLowerCase().includes("safari");
}

const [toggleLoop, handleLoop, checkLoop] = looper();
const [playMuzikk, pauseMuzikk, togglePP] = playPauser(audioPlayerEl);
const stopMuzikk = stopper(audioPlayerEl);
const [shuffle, toggleShuffle] = shuffler(playerState);
const [
  nextMuzikk,
  previousMuzikk,
  removeSongFromPlaylist,
  highlightSongInPlaylist,
] = playlister(playerState);
const [setVolume, mute] = volumer();
const [setPlaybackSpeed] = playebackSpeeder();

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
nextExpandEl?.addEventListener("click", nextMuzikk);
prevEl?.addEventListener("click", previousMuzikk);
prevExpandEl?.addEventListener("click", previousMuzikk);
shuffleEl?.addEventListener("click", toggleShuffle);
shuffleExpandEl?.addEventListener("click", toggleShuffle);
loopEl?.addEventListener("click", toggleLoop);
loopExpandEl?.addEventListener("click", toggleLoop);

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
  if (!!playPauseToggleEl) playPauseToggleEl.disabled = null;
  if (!!playPauseToggleExapndedEl) playPauseToggleExapndedEl.disabled = null;
  if (!!shuffleEl) shuffleEl.disabled = null;
  if (!!shuffleExpandEl) shuffleExpandEl.disabled = null;
  if (!!nextEl) nextEl.disabled = null;
  if (!!nextExpandEl) nextExpandEl.disabled = null;
  if (!!prevEl) prevEl.disabled = null;
  if (!!prevExpandEl) prevExpandEl.disabled = null;
  if (!!loopEl) loopEl.disabled = null;
  if (!!loopExpandEl) loopExpandEl.disabled = null;

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

let playerStartY = 0;

playerEl?.addEventListener(
  "touchstart",
  (e) => {
    playerStartY = e.touches[0].pageY;
  },
  { passive: true },
);

playerEl?.addEventListener(
  "touchmove",
  async (e) => {
    const y = e.touches[0].pageY;
    if (y > playerStartY + 75) {
      collapse();
    }
    if (y < playerStartY - 25) {
      expand();
    }
  },
  { passive: true },
);

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
window.Player.downloadPlaylistToDevice = downloadPlaylistToDevice;
window.Player.showPlayer = show;
window.Player.hidePlayer = hide;
window.Player.playSingleSong = playSingleSong;
window.Player.playSingleSongId = playSingleSongId;
window.Player.playSingleSongNext = playSingleSongNext;
window.Player.playSingleSongNextId = playSingleSongNextId;
window.Player.playSongFromPlaylist = playSongFromPlaylist;
window.Player.playSongFromPlaylistId = playSongFromPlaylistId;
window.Player.playPlaylistNext = playPlaylistNext;
window.Player.playPlaylistNextId = playPlaylistNextId;
window.Player.removeSongFromPlaylist = removeSongFromPlaylist;
window.Player.addSongToQueue = appendSongToCurrentQueue;
window.Player.addSongToQueueId = appendSongToCurrentQueueId;
window.Player.appendPlaylistToCurrentQueue = appendPlaylistToCurrentQueue;
window.Player.appendPlaylistToCurrentQueueId = appendPlaylistToCurrentQueueId;
window.Player.stopMuzikk = stopMuzikk;
window.Player.setPlaybackSpeed = setPlaybackSpeed;
window.Player.expand = () => expand();
window.Player.collapse = () => collapse();
