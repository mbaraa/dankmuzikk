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
 * @property {string} media_url
 */

function playMuzikk() {
  audioPlayerEl.muted = null;
  audioPlayerEl.play();
  PlayerUI.setPauseIcon();
}

function pauseMuzikk() {
  audioPlayerEl.pause();
  PlayerUI.setPlayIcon();
}

function stopMuzikk() {
  pauseMuzikk();
  audioPlayerEl.currentTime = 0;
}

function playPauseMuzikk() {
  if (audioPlayerEl.paused) playMuzikk();
  else pauseMuzikk();
}

/**
 * @param {string} songPublicId
 * @param {string} playlistPublicId
 */
async function removeSongFromPlaylist(songPublicId, playlistPublicId) {
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
 * @param {Song} song
 */
function playSong(song) {
  audioPlayerEl.muted = null;
  audioPlayerEl.pause();
  audioPlayerEl.currentTime = 0;
  if (audioPlayerEl.childNodes.length > 0) {
    audioPlayerEl.removeChild(audioPlayerEl.childNodes.item(0));
  }
  const src = document.createElement("source");
  src.setAttribute("type", "audio/mpeg");
  src.setAttribute("src", song.media_url);
  src.setAttribute("preload", "metadata");
  audioPlayerEl.appendChild(src);
  audioPlayerEl.load();

  PlayerUI.setSongName(song.title);
  PlayerUI.setArtistName(song.artist);
  PlayerUI.setSongThumbnail(song.thumbnail_url);
  PlayerUI.setLoadingOff();
  audioPlayerEl.play();
  setMediaSessionMetadata(song);
  PlayerUI.triggerFetchLyrics();
}

/**
 * @param {string} songPublicId
 * @param {string} playlistPublicId
 * @param {"single" | "playlist" | "favorites" | "queue"} from
 */
async function fetchAndPlaySong(songPublicId, playlistPublicId, from) {
  PlayerUI.setLoadingOn();

  let url;
  switch (from) {
    case "single":
      url = "/api/song/play?id=" + songPublicId;
      break;
    case "playlist":
      url = `/api/song/play/playlist?id=${songPublicId}&playlist-id=${playlistPublicId}`;
      break;
    case "queue":
      url = "/api/song/play/queue?id=" + songPublicId;
      break;
    case "favorites":
      url = "/api/song/play/favorites?id=" + songPublicId;
      break;
  }

  let resp = null;

  for (let i = 0; i < 30; i++) {
    resp = await fetch(url, { method: "PUT" })
      .then((res) => res.json())
      .catch((e) => {
        console.error(e);
        return null;
      });
    if (!!resp.media_url) {
      break;
    }
    await Utils.sleep(1000);
  }
  PlayerUI.setLoadingOff();
  if (!resp) {
    alert("Something went wrong when loading the song...");
  }

  playSong(resp);
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
  Utils.showLoading();

  try {
    let song = null;
    let songMetadataFetched = false;

    for (let i = 0; i < 30; i++) {
      const resp = await fetch(`/api/song/play?id=${songPublicId}`, {
        method: "PUT",
      });
      if (!resp.ok) {
        throw new Error("Something went wrong when fetching song's metadata");
      }
      song = await resp.json();
      songMetadataFetched = true;

      if (!!song.media_url) {
        break;
      }
      await Utils.sleep(1000);
    }

    if (!songMetadataFetched) {
      throw new Error("Failed to fetch song metadata after multiple attempts.");
    }

    if (!song.media_url) {
      console.warn("Song not fully downloaded after repeated checks.");
    }

    if (song.media_url) {
      const a = document.createElement("a");
      a.href = song.media_url.replace("muzikkx", "muzikkx-raw");
      a.click();
      a.remove();
      return { ok: true, ...song };
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

/**
 * @param {string} playlistPublicId
 * @param {string} playlistTitle
 */
async function downloadPlaylistToDevice(playlistPublicId, playlistTitle) {
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

function handleUIEvents() {
  {
    const __handler = (e) => {
      e.stopImmediatePropagation();
      e.preventDefault();
      const seekTime = Number(e.target.value);
      audioPlayerEl.currentTime = seekTime;
      setManualSeekOn();
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
    }
  }
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
    setManualSeekOn();
    let seekTo = -10;
    if (audioPlayerEl.currentTime + seekTo < 0) {
      seekTo = 0;
    }
    audioPlayerEl.currentTime += seekTo;
  });

  navigator.mediaSession.setActionHandler("seekforward", () => {
    setManualSeekOn();
    let seekTo = +10;
    if (audioPlayerEl.currentTime + seekTo > audioPlayerEl.duration) {
      seekTo = 0;
    }
    audioPlayerEl.currentTime += seekTo;
  });

  navigator.mediaSession.setActionHandler("seekto", (a) => {
    setManualSeekOn();
    const seekTime = Number(a.seekTime);
    audioPlayerEl.currentTime = seekTime;
  });

  navigator.mediaSession.setActionHandler("previoustrack", async () => {
    PlayerUI.__elements.prevEl.click();
  });

  navigator.mediaSession.setActionHandler("nexttrack", async () => {
    PlayerUI.__elements.nextEl.click();
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

function lyricsParter() {
  let partsMs = [];

  /**
   * @param {string} ms
   */
  const __getPart = (ms) => {
    let part = "";
    for (let i = partsMs.length - 1; i >= 0; i--) {
      if (ms.localeCompare(partsMs[i]) > 0) {
        part = partsMs[i];
        break;
      }
    }

    return part;
  };

  const __getParts = () => {
    return partsMs;
  };

  const __setParts = (parts) => {
    partsMs = parts;
    resetLyricsHighlighter();
  };

  return [__getPart, __getParts, __setParts];
}

function lyricsHighlighter() {
  let lastLyricsPartMs = "";
  let manualSeek = false;

  const __handle = (event) => {
    const ms = getLyricsPartMs(Utils.formatTimeMs(event.target.currentTime));
    if (!ms) {
      return;
    }
    if (ms.localeCompare(lastLyricsPartMs) <= 0 && !manualSeek) return;
    manualSeek = false;

    const lastPart = document.getElementById("lyrics-part-" + lastLyricsPartMs);
    if (!!lastPart) {
      lastPart.style.color = "var(--secondary-color-69)";
      lastPart.style.fontSize = "medium";
      lastPart.style.fontWeight = "lighter";
    }

    const currentPart = document.getElementById("lyrics-part-" + ms);
    if (!!currentPart) {
      currentPart.style.color = "var(--secondary-color)";
      currentPart.style.fontSize = "x-large";
      currentPart.style.fontWeight = "normal";
      currentPart.scrollIntoView({
        behavior: "smooth",
        block: "center",
        inline: "nearest",
      });
    }

    lastLyricsPartMs = ms;
  };

  const __setManualSeekOn = () => {
    manualSeek = true;
  };

  const __reset = () => {
    manualSeek = false;
    lastLyricsPartMs = "";
  };

  return [__handle, __setManualSeekOn, __reset];
}

handleUIEvents();
handleMediaSessionEvents();

const [setVolume, mute] = volumer();
const setPlaybackSpeed = playebackSpeeder();
const [getLyricsPartMs, getLyricsPartsMs, setLyicsPartsMs] = lyricsParter();
const [handleLyricsHighlight, setManualSeekOn, resetLyricsHighlighter] =
  lyricsHighlighter();

window.Player = {
  playPauseMuzikk,
  stopMuzikk,
  downloadSongToDevice,
  downloadPlaylistToDevice,
  removeSongFromPlaylist,
  setPlaybackSpeed,
  playSong,
  fetchAndPlaySong,
  setLyicsPartsMs,
  handleLyricsHighlight,
};
