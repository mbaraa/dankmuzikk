"use strict";

const shuffleEl1 = document.getElementById("shuffle"),
  nextEl1 = document.getElementById("next"),
  prevEl1 = document.getElementById("prev");

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
  #playlist;
  #currentSongIndex;

  /**
   * @param {Playlist} playlist
   */
  constructor(playlist) {
    this.#playlist = playlist;
    this.#currentSongIndex = 0;
  }

  /**
   * @param {string} songYtId
   */
  play(songYtId = "") {
    this.setSongNotPlayingStyle();
    this.#currentSongIndex = this.#playlist.songs.findIndex(
      (song) => song.yt_id === songYtId,
    );
    if (this.#currentSongIndex < 0) {
      this.#currentSongIndex = 0;
    }
    const songToPlay = this.#playlist.songs[this.#currentSongIndex];
    Player.playSong(songToPlay, true);
    nextEl1.style.display = "block";
    prevEl1.style.display = "block";
    shuffleEl1.style.display = "block";
    this.#updateSongPlays();
    this.setSongPlayingStyle();
  }

  next(shuffle = false, loop = false) {
    this.setSongNotPlayingStyle();
    if (
      !loop &&
      !shuffle &&
      this.#currentSongIndex + 1 >= this.#playlist.songs.length
    ) {
      Player.stopMuzikk();
      return;
    }
    this.#currentSongIndex = shuffle
      ? Math.floor(Math.random() * this.#playlist.songs.length)
      : loop && this.#currentSongIndex + 1 >= this.#playlist.songs.length
        ? 0
        : this.#currentSongIndex + 1;
    const songToPlay = this.#playlist.songs[this.#currentSongIndex];
    Player.playSong(songToPlay, true);
    this.#updateSongPlays();
    this.setSongPlayingStyle();
  }

  previous(shuffle = false, loop = false) {
    this.setSongNotPlayingStyle();
    if (!loop && !shuffle && this.#currentSongIndex - 1 < 0) {
      Player.stopMuzikk();
      return;
    }
    this.#currentSongIndex = shuffle
      ? Math.floor(Math.random() * this.#playlist.songs.length)
      : loop && this.#currentSongIndex - 1 < 0
        ? this.#playlist.songs.length - 1
        : this.#currentSongIndex - 1;
    this.setSongNotPlayingStyle();
    const songToPlay = this.#playlist.songs[this.#currentSongIndex];
    Player.playSong(songToPlay, true);
    this.#updateSongPlays();
    this.setSongPlayingStyle();
  }

  removeSong(songYtId) {
    const songIndex = this.#playlist.songs.findIndex(
      (song) => song.yt_id === songYtId,
    );
    if (songIndex < 0) {
      return;
    }
    this.#playlist.songs.splice(songIndex, 1);
  }

  setSongPlayingStyle() {
    const songEl = document.getElementById(
      "song-" + this.#playlist.songs[this.#currentSongIndex].yt_id,
    );
    if (!songEl) {
      return;
    }
    songEl.style.backgroundColor = "var(--accent-color-30)";
    songEl.scrollIntoView();
  }

  setSongNotPlayingStyle() {
    for (const song of this.#playlist.songs) {
      const songEl = document.getElementById("song-" + song.yt_id);
      if (!songEl) {
        return;
      }
      songEl.style.backgroundColor = "var(--secondary-color-20)";
    }
  }

  async #updateSongPlays() {
    await fetch(
      "/api/increment-song-plays?" +
        new URLSearchParams({
          "song-id": this.#playlist.songs[this.#currentSongIndex].yt_id,
          "playlist-id": this.#playlist.public_id,
        }).toString(),
      {
        method: "PUT",
      },
    ).catch((err) => console.error(err));
  }
}

/**
 * @param {string} songYtId
 */
function removeSongFromPlaylist(songYtId) {
  if (!Player.playlistPlayer) {
    return;
  }
  Player.playlistPlayer.removeSong(songYtId);
}

/**
 * @param {Playlist} playlist
 */
function playPlaylist(playlist) {
  Player.playlistPlayer = new PlaylistPlayer(playlist);
  Player.playlistPlayer.play();
}

/**
 * @param {string} songId
 * @param {Playlist} playlist
 */
function playSongFromPlaylist(songId, playlist) {
  Player.playlistPlayer = new PlaylistPlayer(playlist);
  Player.playlistPlayer.play(songId);
}

if (!window.Player) {
  window.Player = {};
}

Player.playlistPlayer = null;
window.Player.playPlaylist = playPlaylist;
window.Player.playSongFromPlaylist = playSongFromPlaylist;
window.Player.removeSongFromPlaylist = removeSongFromPlaylist;
