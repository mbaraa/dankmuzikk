"use strict";

const playerButtonsIcons = {
  play: `<img class="w-[50px] h-[50px]" src="/static/icons/play-icon.svg" alt="Play"/>`,
  pause: `<img class="w-[50px] h-[50px]" src="/static/icons/pause-icon.svg" alt="Pause"/>`,
  loop: `<img class="w-[40px]" src="/static/icons/loop-icon.svg" alt="Loop"/>`,
  loopOnce: `<img class="w-[40px]" src="/static/icons/loop-once-icon.svg" alt="Loop Once"/>`,
  loopOff: `<img class="w-[40px]" src="/static/icons/loop-off-icon.svg" alt="Loop Off"/>`,
  shuffle: `<img src="/static/icons/shuffle-icon.svg" alt="Shuffle"/>`,
  shuffleOff: `<img src="/static/icons/shuffle-off-icon.svg" alt="Shuffle"/>`,
  loading: `<div class="loader !h-10 !w-10"></div>`,
};

const loopModes = [
  { icon: "loop-off-icon.svg", mode: "OFF" },
  { icon: "loop-once-icon.svg", mode: "ONCE" },
  { icon: "loop-icon.svg", mode: "ALL" },
];

const playPauseToggleEl = document.getElementById("play"),
  shuffleEl = document.getElementById("shuffle"),
  nextEl = document.getElementById("next"),
  prevEl = document.getElementById("prev"),
  loopEl = document.getElementById("loop"),
  songNameEl = document.getElementById("song-name"),
  artistNameEl = document.getElementById("artist-name"),
  songSeekBarEl = document.getElementById("song-seek-bar"),
  songDurationEl = document.getElementById("song-duration"),
  songCurrentTimeEl = document.getElementById("song-current-time"),
  songImageEl = document.getElementById("song-image"),
  audioPlayerEl = document.getElementById("audio-player"),
  muzikkContainerEl = document.getElementById("muzikk");

let shuffleSongs = false;
let currentLoopIdx = 0;
/**
 * @type{PlaylistPlayer}
 */
let currentPlaylistPlayer;

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

class PlaylistPlayer {
  #currentPlaylist;
  #currentSongIndex;

  /**
   * @param {Playlist} playlist
   */
  constructor(playlist) {
    this.#currentPlaylist = playlist;
    this.#currentSongIndex = 0;
  }

  /**
   * @param {string} songYtIt
   */
  play(songYtIt = "") {
    this.#currentSongIndex = this.#currentPlaylist.songs.findIndex(
      (song) => song.yt_id === songYtIt,
    );
    if (this.#currentSongIndex < 0) {
      this.#currentSongIndex = 0;
    }
    const songToPlay = this.#currentPlaylist.songs[this.#currentSongIndex];
    playNewSong(
      songToPlay.yt_id,
      songToPlay.thumbnail_url,
      songToPlay.title,
      songToPlay.artist,
      songToPlay.duration,
    );
    this.#updateSongPlays();
  }

  next(shuffle = false, loop = false) {
    if (
      !loop &&
      !shuffle &&
      this.#currentSongIndex + 1 >= this.#currentPlaylist.songs.length
    ) {
      stopMuzikk();
      return;
    }

    this.#currentSongIndex = shuffle
      ? Math.floor(Math.random() * this.#currentPlaylist.songs.length)
      : loop && this.#currentSongIndex + 1 >= this.#currentPlaylist.songs.length
        ? 0
        : this.#currentSongIndex + 1;

    const songToPlay = this.#currentPlaylist.songs[this.#currentSongIndex];
    playNewSong(
      songToPlay.yt_id,
      songToPlay.thumbnail_url,
      songToPlay.title,
      songToPlay.artist,
      songToPlay.duration,
    );
    this.#updateSongPlays();
  }

  previous(shuffle = false, loop = false) {
    if (!loop && !shuffle && this.#currentSongIndex - 1 < 0) {
      stopMuzikk();
      return;
    }
    this.#currentSongIndex = shuffle
      ? Math.floor(Math.random() * this.#currentPlaylist.songs.length)
      : loop && this.#currentSongIndex - 1 < 0
        ? this.#currentPlaylist.songs.length - 1
        : this.#currentSongIndex - 1;
    const songToPlay = this.#currentPlaylist.songs[this.#currentSongIndex];
    playNewSong(
      songToPlay.yt_id,
      songToPlay.thumbnail_url,
      songToPlay.title,
      songToPlay.artist,
      songToPlay.duration,
    );
    this.#updateSongPlays();
  }

  /**
   * nextSong returns the next song if applicable
   * @returns{Song}
   */
  nextSong(shuffle = false, loop = false) {
    if (
      !loop &&
      !shuffle &&
      this.#currentSongIndex + 1 >= this.#currentPlaylist.songs.length
    ) {
      stopMuzikk();
      return;
    }

    const nextIndex = shuffle
      ? Math.floor(Math.random() * this.#currentPlaylist.songs.length)
      : loop && this.#currentSongIndex + 1 >= this.#currentPlaylist.songs.length
        ? 0
        : this.#currentSongIndex + 1;

    return this.#currentPlaylist.songs[nextIndex];
  }

  async #updateSongPlays() {
    await fetch(
      "/api/increment-song-plays?" +
        new URLSearchParams({
          "song-id": this.#currentPlaylist.songs[this.#currentSongIndex].yt_id,
          "playlist-id": this.#currentPlaylist.public_id,
        }).toString(),
      {
        method: "PUT",
      },
    ).catch((err) => console.error(err));
  }
}

/**
 * @param {Playlist} playlist
 */
function playPlaylist(playlist) {
  currentPlaylistPlayer = new PlaylistPlayer(playlist);
  currentPlaylistPlayer.play();
}

/**
 * @param {string} songId
 * @param {Playlist} playlist
 */
function playSongFromPlaylist(songId, playlist) {
  currentPlaylistPlayer = new PlaylistPlayer(playlist);
  currentPlaylistPlayer.play(songId);
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
  if (!currentPlaylistPlayer) {
    return;
  }
  currentPlaylistPlayer.next(
    shuffleSongs,
    loopModes[currentLoopIdx].mode === "ALL",
  );
}

function previousMuzikk() {
  if (!currentPlaylistPlayer) {
    return;
  }
  currentPlaylistPlayer.previous(
    shuffleSongs,
    loopModes[currentLoopIdx].mode === "ALL",
  );
}

function playMuzikk() {
  audioPlayerEl.play();
  playPauseToggleEl.innerHTML = playerButtonsIcons.pause;
}

function pauseMuzikk() {
  audioPlayerEl.pause();
  playPauseToggleEl.innerHTML = playerButtonsIcons.play;
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
  if (!currentPlaylistPlayer) {
    window.alert("Shuffling can't be enabled for a single song!");
    return;
  }
  if (shuffleSongs) {
    shuffleEl.innerHTML = playerButtonsIcons.shuffleOff;
  } else {
    shuffleEl.innerHTML = playerButtonsIcons.shuffle;
  }
  shuffleSongs = !shuffleSongs;
}

/**
 * @param {Song} song
 */
async function downloadSong(song) {
  return await fetch(
    "/api/song/download?" + new URLSearchParams(song).toString(),
  ).catch((err) => console.error(err));
}

/**
 * @param {Song} song
 */
async function playNewSong(song) {
  playPauseToggleEl.innerHTML = playerButtonsIcons.loading;
  document.body.style.cursor = "progress";
  Utils.showLoading();

  await downloadSong(song).then(() => {
    stopMuzikk();
    audioPlayerEl.src = `/music/${song.yt_id}.mp3`;
    audioPlayerEl.load();
  });

  setMediaSession(song);
  showPlayer();

  if (song.title) {
    songNameEl.innerHTML = song.title;
    songNameEl.title = song.title;
    if (song.title.length > Utils.getTextWidth()) {
      songNameEl.parentElement.classList.add("marquee");
    } else {
      songNameEl.parentElement.classList.remove("marquee");
    }
  }
  if (song.artist) {
    artistNameEl.innerHTML = song.artist;
    artistNameEl.title = song.artist;
    if (song.artist.length > Utils.getTextWidth()) {
      artistNameEl.parentElement.classList.add("marquee");
    } else {
      artistNameEl.parentElement.classList.remove("marquee");
    }
  }

  playMuzikk();
  songImageEl.style.backgroundImage = `url("${song.thumbnail_url}")`;
}

function showPlayer() {
  muzikkContainerEl.style.display = "block";
}

function hidePlayer() {
  muzikkContainerEl.style.display = "none";
  audioPlayerEl.stopMuzikk();
}

playPauseToggleEl.addEventListener("click", playPauseToggle);
nextEl.addEventListener("click", nexMuzikk);
prevEl.addEventListener("click", previousMuzikk);
shuffleEl.addEventListener("click", toggleShuffle);

loopEl.addEventListener("click", (event) => {
  currentLoopIdx = (currentLoopIdx + 1) % loopModes.length;
  event.target.src = "/static/icons/" + loopModes[currentLoopIdx].icon;
});

songSeekBarEl.addEventListener("change", (event) => {
  const seekTime = Number(event.target.value);
  audioPlayerEl.currentTime = seekTime;
});

audioPlayerEl.addEventListener("loadeddata", (event) => {
  playPauseToggleEl.disabled = null;
  shuffleEl.disabled = null;
  nextEl.disabled = null;
  prevEl.disabled = null;
  loopEl.disabled = null;

  const duration = Math.ceil(
    Number.isNaN(event.target.duration) ? 0 : event.target.duration,
  );
  songSeekBarEl.max = Math.ceil(duration);
  songSeekBarEl.value = 0;
  if (songDurationEl) {
    songDurationEl.innerHTML = Utils.formatTime(duration);
  }

  playPauseToggleEl.innerHTML = playerButtonsIcons.pause;
  document.body.style.cursor = "auto";
  Utils.hideLoading();
});

audioPlayerEl.addEventListener("timeupdate", (event) => {
  const currentTime = Math.floor(event.target.currentTime);
  if (songCurrentTimeEl) {
    songCurrentTimeEl.innerHTML = Utils.formatTime(currentTime);
  }
  if (songSeekBarEl) {
    songSeekBarEl.value = Math.ceil(currentTime);
  }

  if (
    !!currentPlaylistPlayer &&
    event.target.currentTime >= event.target.duration * 0.9
  ) {
    downloadSong(
      currentPlaylistPlayer.nextSong(
        shuffleSongs,
        loopModes[currentLoopIdx].mode === "ALL",
      ),
    );
  }
});

audioPlayerEl.addEventListener("ended", () => {
  switch (loopModes[currentLoopIdx].mode) {
    case "OFF":
      if (currentPlaylistPlayer) {
        currentPlaylistPlayer.next(shuffleSongs, false);
        return;
      }
      stopMuzikk();
      break;
    case "ONCE":
      stopMuzikk();
      playMuzikk();
      break;
    case "ALL":
      if (currentPlaylistPlayer) {
        currentPlaylistPlayer.next(shuffleSongs, true);
        return;
      }
      stopMuzikk();
      playMuzikk();
      break;
  }
});

audioPlayerEl.addEventListener("progress", () => {
  console.log("downloading...");
});

window.Player = {
  showPlayer,
  hidePlayer,
  playPlaylist,
  playSongFromPlaylist,
  playNewSong,
};
