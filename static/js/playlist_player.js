"use strict";

const shuffleEl1 = document.getElementById("shuffle"),
  nextEl1 = document.getElementById("next"),
  prevEl1 = document.getElementById("prev");

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
    this.setSongNotPlayingStyle();
    this.#currentSongIndex = this.#currentPlaylist.songs.findIndex(
      (song) => song.yt_id === songYtIt,
    );
    if (this.#currentSongIndex < 0) {
      this.#currentSongIndex = 0;
    }
    const songToPlay = this.#currentPlaylist.songs[this.#currentSongIndex];
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
      this.#currentSongIndex + 1 >= this.#currentPlaylist.songs.length
    ) {
      Player.stopMuzikk();
      return;
    }
    this.#currentSongIndex = shuffle
      ? Math.floor(Math.random() * this.#currentPlaylist.songs.length)
      : loop && this.#currentSongIndex + 1 >= this.#currentPlaylist.songs.length
        ? 0
        : this.#currentSongIndex + 1;
    const songToPlay = this.#currentPlaylist.songs[this.#currentSongIndex];
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
      ? Math.floor(Math.random() * this.#currentPlaylist.songs.length)
      : loop && this.#currentSongIndex - 1 < 0
        ? this.#currentPlaylist.songs.length - 1
        : this.#currentSongIndex - 1;
    this.setSongNotPlayingStyle();
    const songToPlay = this.#currentPlaylist.songs[this.#currentSongIndex];
    Player.playSong(songToPlay, true);
    this.#updateSongPlays();
    this.setSongPlayingStyle();
  }

  removeSong(songYtId) {
    const songIndex = this.#currentPlaylist.songs.findIndex(
      (song) => song.yt_id === songYtId,
    );
    if (songIndex < 0) {
      return;
    }
    this.#currentPlaylist.songs.splice(songIndex, 1);
  }

  setSongPlayingStyle() {
    const songEl = document.getElementById(
      "song-" + this.#currentPlaylist.songs[this.#currentSongIndex].yt_id,
    );
    songEl.style.backgroundColor = "var(--accent-color-30)";
    songEl.scrollIntoView();
  }

  setSongNotPlayingStyle() {
    for (const song of this.#currentPlaylist.songs) {
      document.getElementById("song-" + song.yt_id).style.backgroundColor =
        "var(--secondary-color-20)";
    }
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
