"use strict";

const audioPlayerEl = document.getElementById("muzikk-player");

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

/**
 * @typedef {object} PlayerState
 * @property {LoopMode} loopMode
 * @property {boolean} shuffled
 * @property {string} shuffledPlaylist
 * @property {Playlist} playlist
 * @property {number} currentSongIdx
 * @property {boolean} lyricsLoaded
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
  lyricsLoaded: false,
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
    if (__check(LOOP_MODES.ONCE)) {
      PlayerUI.setLoopOnce();
    } else if (__check(LOOP_MODES.ALL)) {
      PlayerUI.setLoopAll();
    } else if (__check(LOOP_MODES.OFF)) {
      PlayerUI.setLoopOff();
    } else {
      PlayerUI.setLoopOff();
    }
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
 * @param {HTMLAudioElement} audioEl
 *
 * @returns {[Function, Function, Function]}
 */
function playPauser(audioEl) {
  const __play = () => {
    audioEl.muted = null;
    audioEl.play();
    PlayerUI.highlightSong(
      playerState.playlist.songs[playerState.currentSongIdx].public_id,
    );
    PlayerUI.setPauseIcon();
  };
  const __pause = () => {
    audioEl.pause();
    PlayerUI.setPlayIcon();
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
    PlayerUI.unHighlightSong(
      playerState.playlist.songs[playerState.currentSongIdx].public_id,
    );
    PlayerUI.setPlayIcon();
  };
}

/**
 * @param {PlayerState} state
 *
 * @returns {Function}
 */
function shuffler(state) {
  // using Fisher–Yates shuffling algorithm
  // https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
  const __shuffleArray = (a) => {
    let currIdx = a.length;
    while (currIdx != 0) {
      let randIdx = Math.floor(Math.random() * currIdx);
      currIdx--;
      [a[currIdx], a[randIdx]] = [a[randIdx], a[currIdx]];
    }
  };

  /**
   * @param {string} songPublicId
   */
  const __shuffle = (songPublicId) => {
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
      sIdx = state.playlist.songs.findIndex(
        (s) => s.public_id === songPublicId,
      );
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
        if (tmp[i].public_id === tmp[i + 1].public_id) {
          continue;
        }
        state.playlist.songs.push(tmp[i]);
      }
      if (tmp[tmp.length - 1].public_id !== tmp[tmp.length - 2]) {
        state.playlist.songs.push(tmp[tmp.length - 1]);
      }
      const src = audioPlayerEl.childNodes.item(0);
      if (!!src) {
        state.currentSongIdx = state.playlist.songs.findIndex(
          (s) =>
            s.public_id ===
            src.src.substring(src.src.lastIndexOf("/") + 1, src.src.length - 4),
        );
      }
    }
    if (state.shuffled) {
      PlayerUI.setShuffleOn();
    } else {
      PlayerUI.setShuffleOff();
    }
  };

  return [__shuffle, __toggleShuffle];
}

/**
 * @param {PlayerState} state
 *
 * @returns {[Function, Function, Function, Function]}
 */
function playlister(state) {
  const __next = async () => {
    if (checkLoop(LOOP_MODES.ONCE)) {
      stopMuzikk();
      playMuzikk();
      return;
    }
    // check votes to whether repeat the song or not.
    if (
      state.playlist.songs[state.currentSongIdx].votes > 1 &&
      !state.shuffled
    ) {
      const songToPlay = state.playlist.songs[state.currentSongIdx];
      songToPlay.votes--;
      songToPlay.plays++;
      playSongFromPlaylist(songToPlay.public_id, state.playlist);
      PlayerUI.highlightSongInPlaylist(
        songToPlay.public_id,
        state.playlist.songs.map((s) => s.public_id),
      );
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
    PlayerUI.highlightSongInPlaylist(
      songToPlay.public_id,
      state.playlist.songs.map((s) => s.public_id),
    );
    await playSong(songToPlay);
  };

  const __prev = async () => {
    if (checkLoop(LOOP_MODES.ONCE)) {
      stopMuzikk();
      playMuzikk();
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
      playSongFromPlaylist(songToPlay.public_id, state.playlist);
      PlayerUI.highlightSongInPlaylist(
        songToPlay.public_id,
        state.playlist.songs.map((s) => s.public_id),
      );
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
    PlayerUI.highlightSongInPlaylist(
      songToPlay.public_id,
      state.playlist.songs.map((s) => s.public_id),
    );
    await playSong(songToPlay);
  };

  const __remove = (songPublicId, playlistId) => {
    const songIndex = state.playlist.songs.findIndex(
      (song) => song.public_id === songPublicId,
    );
    if (songIndex >= 0) {
      state.playlist.songs.splice(songIndex, 1);
    }

    Utils.showLoading();
    fetch(
      "/api/playlist/song?song-id=" +
        songPublicId +
        "&playlist-id=" +
        playlistId +
        "&remove=true",
      {
        method: "PUT",
      },
    )
      .then((res) => {
        if (res.ok) {
          // TODO: do something about this UI leak
          const songEl = document.getElementById("song-" + songPublicId);
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

  return [__next, __prev, __remove];
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

  return [__setSpeed];
}

/**
 * @param {string} songPublicId
 */
async function downloadSong(songPublicId) {
  try {
    const resp = await fetch("/api/song?id=" + songPublicId).then((res) =>
      res.json(),
    );
    for (let i = 0; i < 30; i++) {
      const song = await fetchSongMeta(songPublicId, false);
      if (song.fully_downloaded) {
        return { ok: true, ...resp };
      }
      await Utils.sleep(1000);
    }
  } catch (err) {
    console.error("oopsie", err);
    return { ok: false };
  }
}

/**
 * @param {string} songPublicId
 * @param {string} songTitle
 */
async function downloadSongToDevice(songPublicId, songTitle) {
  Utils.showLoading();
  await downloadSong(songPublicId)
    .then((song) => {
      const a = document.createElement("a");
      a.href = song.media_url.replace("muzikkx", "muzikkx-raw");
      a.download = `${songTitle}.mp3`;
      a.click();
      a.remove();
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
      return res.json();
    })
    .then((res) => {
      const a = document.createElement("a");
      a.href = res["playlist_download_url"];
      a.download = `${plTitle}.zip`;
      a.click();
      a.remove();
    })
    .finally(() => {
      Utils.hideLoading();
    });
}

/**
 * @param {Song} song
 */
async function playSong(song) {
  playerState.lyricsLoaded = false;

  PlayerUI.setLoadingOn();
  PlayerUI.show();

  let resp = await downloadSong(song.public_id);
  if (!resp.ok) {
    alert("Something went wrong when downloading the song...");
    throw new Error("Something went wrong when downloading the song...");
  }
  stopMuzikk();
  if (audioPlayerEl.childNodes.length > 0) {
    audioPlayerEl.removeChild(audioPlayerEl.childNodes.item(0));
  }
  const src = document.createElement("source");
  src.setAttribute("type", "audio/mpeg");
  src.setAttribute("src", resp.media_url);
  src.setAttribute("preload", "metadata");
  audioPlayerEl.appendChild(src);

  if (isSafari()) {
    setTimeout(80);
  }
  audioPlayerEl.load();

  PlayerUI.setSongName(song.title);
  PlayerUI.setArtistName(song.artist);
  PlayerUI.setSongThumbnail(song.thumbnail_url);
  setMediaSessionMetadata(song);
  playMuzikk();
}

/**
 * @param {string} songPublicId
 *
 * @returns {Promise<Song| never>}
 */
async function fetchSongMeta(songPublicId, displayLoader = true) {
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
  playerState.lyricsLoaded = false;

  playerState.playlist = {
    title: "Queue",
    songs_count: 1,
    public_id: "",
    songs: [song],
  };
  playerState.currentSongIdx = 0;

  await window.Utils.retryer(async () => {
    return await playSong(song);
  });
}

/**
 * @param {string} songPublicId
 */
async function playSingleSongId(songPublicId) {
  try {
    await fetchSongMeta(songPublicId).then(
      async (song) => await playSingleSong(song),
    );
  } catch (err) {
    return err;
  }
}

/**
 * @param {Song} song
 */
async function playSingleSongNext(song) {
  if (playerState.playlist.songs.length === 0) {
    try {
      await playSingleSong(song);
    } catch (err) {
      return err;
    }
  }
  if (!song.public_id) {
    return;
  }
  song.votes = 1;
  playerState.playlist.songs.splice(playerState.currentSongIdx + 1, 0, song);
  alert(`Playing ${song.title} next!`);
}

/**
 * @param {string} songPublicId
 */
async function playSingleSongNextId(songPublicId) {
  playerState.lyricsLoaded = false;
  try {
    const song = await fetchSongMeta(songPublicId);
    await playSingleSongNext(song);
  } catch (err) {
    return err;
  }
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
    playSongFromPlaylist(playlist.songs[0].public_id, playlist);
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
    playSongFromPlaylist(playlist.songs[0].public_id, playlist);
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
 * @param {string} songPublicId
 * @param {Playlist} playlist
 */
async function playSongFromPlaylist(songPublicId, playlist) {
  if (
    playerState.shuffled &&
    playerState.shuffledPlaylist !== playlist.public_id
  ) {
    playerState.playlist = playlist;
    shuffle(songPublicId);
  }
  const songIdx = playlist.songs.findIndex((s) => s.public_id === songPublicId);
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
  PlayerUI.highlightSongInPlaylist(
    songToPlay.public_id,
    playerState.playlist.songs.map((s) => s.public_id),
  );
  await playSong(songToPlay);
}

/**
 * @param {string} songPublicId
 * @param {string} playlistPubId
 */
async function playSongFromPlaylistId(songPublicId, playlistPubId) {
  const playlist = await fetchPlaylistMeta(playlistPubId);
  await playSongFromPlaylist(songPublicId, playlist);
  playerState.lyricsLoaded = false;
}

/**
 * @param {Song} song
 */
async function appendSongToCurrentQueue(song) {
  if (playerState.playlist.songs.length === 0) {
    try {
      await playSingleSong(song);
    } catch (err) {
      return err;
    }
  }
  song.votes = 1;
  playerState.playlist.songs.push(song);
  alert(`Added ${song.title} to the queue!`);
}

/**
 * @param {string} songPublicId
 */
async function appendSongToCurrentQueueId(songPublicId) {
  const song = await fetchSongMeta(songPublicId);
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
  return /^((?!chrome|android).)*safari/i.test(navigator.userAgent);
}

const [toggleLoop, handleLoop, checkLoop] = looper();
const [playMuzikk, pauseMuzikk, togglePP] = playPauser(audioPlayerEl);
const stopMuzikk = stopper(audioPlayerEl);
const [shuffle, toggleShuffle] = shuffler(playerState);
const [nextMuzikk, previousMuzikk, removeSongFromPlaylist] =
  playlister(playerState);
const [setVolume, mute] = volumer();
const [setPlaybackSpeed] = playebackSpeeder();

PlayerUI.__elements.playPauseToggleEl.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  togglePP();
});

PlayerUI.__elements.playPauseToggleExapndedEl?.addEventListener(
  "click",
  (event) => {
    event.stopImmediatePropagation();
    event.preventDefault();
    togglePP();
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

(() => {
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
})();

(() => {
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
})();

function loadSongLyrics() {
  console.log("loading lyrics");
  if (playerState.lyricsLoaded) {
    return;
  }

  document.getElementById("current-song-lyrics").innerHTML = "";
  const songPublicId =
    playerState.playlist.songs[playerState.currentSongIdx].public_id;
  htmx
    .ajax("GET", "/api/song/lyrics?id=" + songPublicId, {
      target: "#current-song-lyrics",
      swap: "innerHTML",
    })
    .then(() => {
      playerState.lyricsLoaded = true;
    })
    .catch(() => {
      alert("Lyrics fetching went berzerk...");
    });
}

audioPlayerEl.addEventListener("loadeddata", (event) => {
  PlayerUI.enableButtons();
  PlayerUI.setSongDuration(event.target.duration);

  PlayerUI.setLoadingOff();
});

audioPlayerEl.addEventListener("timeupdate", (event) => {
  PlayerUI.setSongCurrentTime(event.target.currentTime);
});

audioPlayerEl.addEventListener("ended", () => {
  handleLoop();
});

audioPlayerEl.addEventListener("progress", () => {
  console.log("downloading...");
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
  stopMuzikk,
  setPlaybackSpeed,
  loadSongLyrics,
};
