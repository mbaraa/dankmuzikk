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
 * @param {HTMLElement} el
 * @param {string} icon
 */
const setPlayerButtonIcon = (el, icon) => {
  if (!!el && !!icon) {
    el.innerHTML = icon;
  }
};

function disableButtons() {}

function enableButtons() {
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

function setLoopOnce() {
  setPlayerButtonIcon(loopEl, PlayerIcons.loopOnce);
  setPlayerButtonIcon(loopExpandEl, PlayerIcons.loopOnce);
}

function setLoopAll() {
  setPlayerButtonIcon(loopEl, PlayerIcons.loop);
  setPlayerButtonIcon(loopExpandEl, PlayerIcons.loop);
}

function setLoopOff() {
  setPlayerButtonIcon(loopEl, PlayerIcons.loopOff);
  setPlayerButtonIcon(loopExpandEl, PlayerIcons.loopOff);
}

function setShuffleOn() {
  setPlayerButtonIcon(shuffleEl, PlayerIcons.shuffle);
  setPlayerButtonIcon(shuffleExpandEl, PlayerIcons.shuffle);
}

function setShuffleOff() {
  setPlayerButtonIcon(shuffleEl, PlayerIcons.shuffleOff);
  setPlayerButtonIcon(shuffleExpandEl, PlayerIcons.shuffleOff);
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

function muteVolume() {}

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

/**
 * @param {string} songPublicId
 */
function highlightSong(songPublicId) {
  const songEl = document.getElementById("song-" + songPublicId);
  if (!!songEl) {
    songEl.style.backgroundColor = "var(--accent-color-30)";
    songEl.scrollIntoView();
  }
}

/**
 * @param {string} songPublicId
 */
function unHighlightSong(songPublicId) {
  const songEl = document.getElementById("song-" + songPublicId);
  if (!!songEl) {
    songEl.style.backgroundColor = "#ffffff00";
  }
}

function highlightSongInPlaylist(songPublicId, playlistSongIds) {
  for (const songId of playlistSongIds) {
    if (songPublicId === songId) {
      highlightSong(songId);
    } else {
      unHighlightSong(songId);
    }
  }
}

//
// EVENTS
//

document
  .getElementById("expanded-mobile-player-lyrics")
  ?.addEventListener("touchmove", (event) => {
    event.stopImmediatePropagation();
  });

document
  .getElementById("expanded-mobile-player-queue")
  ?.addEventListener("touchmove", (event) => {
    event.stopImmediatePropagation();
  });

window.PlayerUI = {
  __elements: {
    // collapsed player's elements
    playPauseToggleEl,
    shuffleEl,
    nextEl,
    prevEl,
    loopEl,
    songNameEl,
    artistNameEl,
    songSeekBarEl,
    volumeSeekBarEl,
    songDurationEl,
    songCurrentTimeEl,
    songImageEl,
    playerEl,
    collapsedMobilePlayer,
    // expanded player's elements
    playPauseToggleExapndedEl,
    shuffleExpandEl,
    nextExpandEl,
    prevExpandEl,
    loopExpandEl,
    songNameExpandedEl,
    artistNameExpandedEl,
    songSeekBarExpandedEl,
    volumeSeekBarExpandedEl,
    songDurationExpandedEl,
    songCurrentTimeExpandedEl,
    songImageExpandedEl,
    expandedMobilePlayer,
  },

  disableButtons,
  enableButtons,
  setPlayIcon,
  setPauseIcon,
  setSongName,
  setArtistName,
  setSongThumbnail,
  setSongDuration,
  setSongCurrentTime,
  setLoopOnce,
  setLoopAll,
  setLoopOff,
  setShuffleOn,
  setShuffleOff,
  expandMobilePlayer,
  collapseMobilePlayer,
  setLoadingOn,
  setLoadingOff,
  highlightSong,
  unHighlightSong,
  highlightSongInPlaylist,
  setVolumeLevel,
};
