"use strict";

const playerButtonsIcons = {
  play: `<img style="width: 75px;" src="/static/icons/play-icon.svg" alt="Play"/>`,
  pause: `<img style="width: 75px;" src="/static/icons/pause-icon.svg" alt="Pause"/>`,
  loop: `<img style="width: 30px;" src="/static/icons/loop-icon.svg" alt="Loop"/>`,
  loopOnce: `<img style="width: 30px;" src="/static/icons/loop-once-icon.svg" alt="Loop Once"/>`,
  loopOff: `<img style="width: 30px;" src="/static/icons/loop-off-icon.svg" alt="Loop Off"/>`,
  shuffle: `<img style="width: 30px;" src="/static/icons/shuffle-icon.svg" alt="Shuffle"/>`,
  shuffleOff: `<img style="width: 30px;" src="/static/icons/shuffle-off-icon.svg" alt="Shuffle"/>`,
  loading: `<div style="width: 50px; height: 50px;" class="loader"></div>`,
};

const loopModes = [
  { icon: "loop-off-icon.svg", mode: "OFF" },
  { icon: "loop-once-icon.svg", mode: "ONCE" },
  { icon: "loop-icon.svg", mode: "ALL" },
];

const playPauseToggleEl = document.getElementById("play"),
  playPauseToggleExapndedEl = document.getElementById("play-expand"),
  shuffleEl = document.getElementById("shuffle"),
  nextEl = document.getElementById("next"),
  prevEl = document.getElementById("prev"),
  loopEl = document.getElementById("loop"),
  songNameEl = document.getElementById("song-name"),
  songNameExpandedEl = document.getElementById("song-name-expanded"),
  artistNameEl = document.getElementById("artist-name"),
  artistNameExpandedEl = document.getElementById("artist-name-expanded"),
  songSeekBarEl = document.getElementById("song-seek-bar"),
  songSeekBarExpandedEl = document.getElementById("song-seek-bar-expanded"),
  songDurationEl = document.getElementById("song-duration"),
  songCurrentTimeEl = document.getElementById("song-current-time"),
  songCurrentTimeExpandedEl = document.getElementById(
    "song-current-time-expanded",
  ),
  songImageEl = document.getElementById("song-image"),
  songImageExpandedEl = document.getElementById("song-image-expanded"),
  audioPlayerEl = document.getElementById("audio-player"),
  muzikkContainerEl = document.getElementById("muzikk"),
  zePlayerEl = document.getElementById("ze-player"),
  zeCollapsedMobilePlayer = document.getElementById(
    "ze-collapsed-mobile-player",
  ),
  zeExpandedMobilePlayer = document.getElementById("ze-expanded-mobile-player");

let shuffleSongs = false;
let currentLoopIdx = 0;
/**
 * @type{PlayerUI}
 */
let ui;

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

class PlayerUI {
  show() {
    muzikkContainerEl.style.display = "block";
  }

  hide() {
    muzikkContainerEl.style.display = "none";
  }

  expand() {
    if (!zePlayerEl.classList.contains("exapnded")) {
      zePlayerEl.classList.add("exapnded");
      zeCollapsedMobilePlayer.classList.add("hidden");
      zeExpandedMobilePlayer.classList.remove("hidden");
    }
  }

  collapse() {
    if (zePlayerEl.classList.contains("exapnded")) {
      zePlayerEl.classList.remove("exapnded");
      zeExpandedMobilePlayer.classList.add("hidden");
      zeCollapsedMobilePlayer.classList.remove("hidden");
    }
  }

  /**
   * @param {boolean} isPlay
   */
  setPlay(isPlay = false) {
    if (isPlay) {
      playPauseToggleEl.innerHTML = playerButtonsIcons.pause;
      if (!!playPauseToggleExapndedEl) {
        playPauseToggleExapndedEl.innerHTML = playerButtonsIcons.pause;
      }
    } else {
      playPauseToggleEl.innerHTML = playerButtonsIcons.play;
      if (!!playPauseToggleExapndedEl) {
        playPauseToggleExapndedEl.innerHTML = playerButtonsIcons.play;
      }
    }
  }

  /**
   * @param {boolean} isShuffle
   */
  setShuffle(isShuffle = false) {
    if (isShuffle) {
      shuffleEl.innerHTML = playerButtonsIcons.shuffleOff;
    } else {
      shuffleEl.innerHTML = playerButtonsIcons.shuffle;
    }
  }

  /**
   * @param {boolean} loading
   * @param {string} fallback is used when loading is false, that is to reset
   *     the loading thingy
   */
  setLoading(loading, fallback) {
    if (loading) {
      playPauseToggleEl.innerHTML = playerButtonsIcons.loading;
      if (!!playPauseToggleExapndedEl) {
        playPauseToggleExapndedEl.innerHTML = playerButtonsIcons.loading;
      }
      document.body.style.cursor = "progress";
      return;
    }
    if (fallback) {
      playPauseToggleEl.innerHTML = fallback;
      if (!!playPauseToggleExapndedEl) {
        playPauseToggleExapndedEl.innerHTML = fallback;
      }
    }
    document.body.style.cursor = "auto";
  }

  /**
   * @param {Song} song
   */
  setDetails(song) {
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

  /**
   * @param {number} time
   */
  setCurrentTime(time) {
    const currentTime = Math.floor(time);
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
  }

  /**
   * @param {number} duration
   */
  setDuration(duration) {
    if (isNaN(duration)) {
      duration = 0;
    }
    songSeekBarEl.max = Math.ceil(duration);
    songSeekBarEl.value = 0;
    if (songDurationEl) {
      songDurationEl.innerHTML = Utils.formatTime(duration);
    }
  }
}

/**
 * @param {Song} song
 */
function setMediaSession(song) {
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

  navigator.mediaSession.setActionHandler("play", () => {
    playPauseToggle();
  });

  navigator.mediaSession.setActionHandler("pause", () => {
    playPauseToggle();
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
    nexMuzikk();
  });

  navigator.mediaSession.setActionHandler("stop", () => {
    stopMuzikk();
  });
}

function nexMuzikk() {
  if (!Player.playlistPlayer) {
    return;
  }
  Player.playlistPlayer.next(
    shuffleSongs,
    loopModes[currentLoopIdx].mode === "ALL",
  );
}

function previousMuzikk() {
  if (!Player.playlistPlayer) {
    return;
  }
  Player.playlistPlayer.previous(
    shuffleSongs,
    loopModes[currentLoopIdx].mode === "ALL",
  );
}

function playMuzikk() {
  audioPlayerEl.play();
  ui.setPlay(true);
}

function pauseMuzikk() {
  audioPlayerEl.pause();
  ui.setPlay(false);
}

function stopMuzikk() {
  pauseMuzikk();
  audioPlayerEl.currentTime = 0;
}

function playPauseToggle() {
  if (audioPlayerEl.paused) {
    playMuzikk();
  } else {
    pauseMuzikk();
  }
}

function toggleShuffle() {
  if (!Player.playlistPlayer) {
    alert("Shuffling can't be enabled for a single song!");
    return;
  }
  ui.setShuffle(shuffleSongs);
  shuffleSongs = !shuffleSongs;
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
 * @param {Song} song
 * @param {boolean} inPlaylist
 */
async function playSong(song, inPlaylist) {
  if (!inPlaylist) {
    Player.playlistPlayer = null;
    nextEl.style.display = "none";
    prevEl.style.display = "none";
    shuffleEl.style.display = "none";
    if (currentLoopIdx > 1) {
      currentLoopIdx = 0;
      loopEl.children[0].src =
        "/static/icons/" + loopModes[currentLoopIdx].icon;
    }
  }
  ui.setLoading(true);
  ui.show();

  await downloadSong(song.yt_id).then(() => {
    stopMuzikk();
    audioPlayerEl.src = `/muzikkx/${song.yt_id}.mp3`;
    audioPlayerEl.load();
  });

  ui.setDetails(song);
  setMediaSession(song);
  playMuzikk();
}

function showPlayer() {
  muzikkContainerEl.style.display = "block";
}

function hidePlayer() {
  muzikkContainerEl.style.display = "none";
  audioPlayerEl.stopMuzikk();
}

playPauseToggleEl.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  playPauseToggle();
});

playPauseToggleExapndedEl?.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  playPauseToggle();
});

nextEl?.addEventListener("click", nexMuzikk);
prevEl?.addEventListener("click", previousMuzikk);
shuffleEl?.addEventListener("click", toggleShuffle);

loopEl?.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  if (!Player.playlistPlayer) {
    currentLoopIdx = currentLoopIdx === 0 ? 1 : 0;
  } else {
    currentLoopIdx = (currentLoopIdx + 1) % loopModes.length;
  }
  event.target.src = "/static/icons/" + loopModes[currentLoopIdx].icon;
});

songSeekBarEl.addEventListener("change", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
  const seekTime = Number(event.target.value);
  audioPlayerEl.currentTime = seekTime;
});

songSeekBarEl.addEventListener("click", (event) => {
  event.stopImmediatePropagation();
  event.preventDefault();
});

audioPlayerEl.addEventListener("loadeddata", (event) => {
  playPauseToggleEl.disabled = null;
  if (playPauseToggleExapndedEl) {
    playPauseToggleExapndedEl.disabled = null;
  }
  shuffleEl.disabled = null;
  nextEl.disabled = null;
  prevEl.disabled = null;
  loopEl.disabled = null;

  ui.setDuration(event.target.duration);
  ui.setLoading(false, playerButtonsIcons.pause);
});

audioPlayerEl.addEventListener("timeupdate", (event) => {
  ui.setCurrentTime(event.target.currentTime);
});

audioPlayerEl.addEventListener("ended", () => {
  switch (loopModes[currentLoopIdx].mode) {
    case "OFF":
      stopMuzikk();
      if (Player.playlistPlayer) {
        Player.playlistPlayer.next(shuffleSongs, false);
      }
      break;
    case "ONCE":
      stopMuzikk();
      playMuzikk();
      break;
    case "ALL":
      if (Player.playlistPlayer) {
        Player.playlistPlayer.next(shuffleSongs, true);
        return;
      }
      break;
  }
});

audioPlayerEl.addEventListener("progress", () => {
  console.log("downloading...");
});

document
  .getElementById("collapse-player-button")
  ?.addEventListener("click", (event) => {
    event.stopImmediatePropagation();
    event.preventDefault();
    ui.collapse();
  });

function init() {
  ui = new PlayerUI();
}

init();

window.Player = {};
window.Player.showPlayer = showPlayer;
window.Player.hidePlayer = hidePlayer;
window.Player.playSong = playSong;
window.Player.stopMuzikk = stopMuzikk;
window.Player.expand = () => ui.expand();
window.Player.collapse = () => ui.collapse();
