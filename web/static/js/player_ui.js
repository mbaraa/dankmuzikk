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

  setSongName,
  setArtistName,
  setSongThumbnail,
  setLoopOnce,
  setLoopAll,
  setLoopOff,
  setShuffleOn,
  setShuffleOff,
  highlightSong,
  unHighlightSong,
  highlightSongInPlaylist,
  setVolumeLevel,
};
