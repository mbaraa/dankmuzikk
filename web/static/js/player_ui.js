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
  playerEl = document.getElementById("ze-player"),
  collapsedMobilePlayer = document.getElementById("ze-collapsed-mobile-player");

// expanded player's elements
const playPauseToggleExapndedEl = document.getElementById("play-expand"),
  nextExpandEl = document.getElementById("next-expand"),
  prevExpandEl = document.getElementById("prev-expand"),
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
 * @param {boolean} loading
 */
function setLoading(loading) {
  if (loading) {
    setPlayerButtonIcon(playPauseToggleEl, PlayerIcons.loading);
    setPlayerButtonIcon(playPauseToggleExapndedEl, PlayerIcons.loading);
    document.body.style.cursor = "progress";
    return;
  }
  setPlayerButtonIcon(playPauseToggleEl, PlayerIcons.pause);
  setPlayerButtonIcon(playPauseToggleExapndedEl, PlayerIcons.pause);
  document.body.style.cursor = "auto";
}

function setLoadingOn() {
  setLoading(true);
}

function setLoadingOff() {
  setLoading(false);
}

function expandMobilePlayer() {
  if (!playerEl.classList.contains("exapnded")) {
    playerEl.classList.add("exapnded");
    collapsedMobilePlayer.classList.add("hidden");
    expandedMobilePlayer.classList.remove("hidden");
  }
}

function collapseMobilePlayer() {
  if (playerEl.classList.contains("exapnded")) {
    playerEl.classList.remove("exapnded");
    collapsedMobilePlayer.classList.remove("hidden");
    expandedMobilePlayer.classList.add("hidden");
  }
}

function enableButtons() {
  if (!!playPauseToggleEl) playPauseToggleEl.disabled = null;
  if (!!playPauseToggleExapndedEl) playPauseToggleExapndedEl.disabled = null;
  if (!!shuffleEl) shuffleEl.disabled = null;
  if (!!nextEl) nextEl.disabled = null;
  if (!!nextExpandEl) nextExpandEl.disabled = null;
  if (!!prevEl) prevEl.disabled = null;
  if (!!prevExpandEl) prevExpandEl.disabled = null;
  if (!!loopEl) loopEl.disabled = null;
}

function setPlayIcon() {
  setPlayerButtonIcon(playPauseToggleEl, PlayerIcons.play);
  setPlayerButtonIcon(playPauseToggleExapndedEl, PlayerIcons.play);
}

function setPauseIcon() {
  setPlayerButtonIcon(playPauseToggleEl, PlayerIcons.pause);
  setPlayerButtonIcon(playPauseToggleExapndedEl, PlayerIcons.pause);
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
 * @param {string} name
 */
function setSongName(name) {
  if (name) {
    songNameEl.innerHTML = name;
    songNameEl.title = name;
    if (name.length > Utils.getTextWidth()) {
      songNameEl.parentElement.classList.add("marquee");
    } else {
      songNameEl.parentElement.classList.remove("marquee");
    }

    if (songNameExpandedEl) {
      songNameExpandedEl.innerHTML = name;
      songNameExpandedEl.title = name;
      if (name.length > Utils.getTextWidth()) {
        songNameExpandedEl.parentElement.classList.add("marquee");
      } else {
        songNameExpandedEl.parentElement.classList.remove("marquee");
      }
    }
  }
}

/**
 * @param {string} name
 */
function setArtistName(name) {
  if (name) {
    if (!!artistNameEl) {
      artistNameEl.innerHTML = name;
      artistNameEl.title = name;
    }

    if (artistNameExpandedEl) {
      artistNameExpandedEl.innerHTML = name;
      artistNameExpandedEl.title = name;
    }
  }
}

function setSongThumbnail(thumbUrl) {
  songImageEl.style.backgroundImage = `url("${thumbUrl}")`;
  songImageEl.innerHTML = "";

  if (songImageExpandedEl) {
    songImageExpandedEl.style.backgroundImage = `url("${thumbUrl}")`;
    songImageExpandedEl.innerHTML = "";
  }
}

function setSongDuration(duration) {
  if (Number.isNaN(duration) || !Number.isFinite(duration)) {
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

function setSongCurrentTime(time) {
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
 * @param {number} level
 */
function setVolumeLevel(level) {
  if (volumeSeekBarEl) {
    volumeSeekBarEl.value = Math.floor(level * 100);
  }
  if (volumeSeekBarExpandedEl) {
    volumeSeekBarExpandedEl.value = Math.floor(level * 100);
  }
}

window.PlayerUI = {
  __elements: {
    nextEl,
    prevEl,
    songSeekBarEl,
    songSeekBarExpandedEl,
    volumeSeekBarExpandedEl,
    volumeSeekBarExpandedEl,
  },

  enableButtons,
  setSongName,
  setArtistName,
  setSongThumbnail,
  setSongDuration,
  setSongCurrentTime,
  expandMobilePlayer,
  collapseMobilePlayer,
  setPlayIcon,
  setPauseIcon,
  setLoadingOn,
  setLoadingOff,
  setVolumeLevel,
};
