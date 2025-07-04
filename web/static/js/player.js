"use strict";

const audioPlayerEl = document.getElementById("muzikk-player");

let setVolume, mute;
let setPlaybackSpeed;

async function init() {
  handleUIEvents();
  handleMediaSessionEvents();

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
function playSongV2(song) {
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
  setMediaSessionMetadata(song);
  audioPlayerEl.play();
}

/**
 * @param {string} songPublicId
 * @param {string} playlistPublicId
 */
async function fetchAndPlaySong(songPublicId, playlistPublicId) {
  setPlayerLoadingOn();

  let resp = null;

  for (let i = 0; i < 30; i++) {
    resp = await fetch(
      "/api/song/play?id=" +
        songPublicId +
        (!!playlistPublicId ? `&playlist-id=${playlistPublicId}` : ""),
      { method: "PUT" },
    )
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
  if (!resp) {
    alert("Something went wrong when loading the song...");
  }

  playSongV2(resp);
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

function togglePlayPause() {
  playPauseMuzikk();
}

function handleUIEvents() {
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
  removeSongFromPlaylist,
  setPlaybackSpeed,
  playSongV2,
  fetchAndPlaySong,
};
