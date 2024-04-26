window.playPauseIcons = {
  playDisabled: `<img src="/static/images/play-disabled-icon.svg" alt="Play"/>`,
  play: `<img src="/static/images/play-icon.svg" alt="Play"/>`,
  pauseDisabled: `<img src="/static/images/pause-disabled-icon.svg" alt="Pause"/>`,
  pause: `<img src="/static/images/pause-icon.svg" alt="Pause"/>`,
  loading: `<div class="loader !h-10 !w-10"></div>`,
};

const playPauseToggleEl = document.getElementById("play");
const shuffleEl = document.getElementById("shuffle");
const nextEl = document.getElementById("next");
const prevEl = document.getElementById("prev");
const songNameEl = document.getElementById("song-name");
const artistNameEl = document.getElementById("artist-name");
const songSeekBarEl = document.getElementById("song-seek-bar");
const songDurationEl = document.getElementById("song-duration");
const songCurrentTimeEl = document.getElementById("song-current-time");
const songImageEl = document.getElementById("song-image");

window.isPlayerReady = false;

function onYouTubeIframeAPIReady() {
  window.player = new YT.Player("yt-player", {
    videoId: "NkE3t4ZrUGQ",
    playerVars: {
      playsinline: 1,
    },
    events: {
      onReady: onPlayerReady,
      onStateChange: onPlayerStateChange,
    },
  });
}

let timeUpdateInterval;

function onPlayerReady(event) {
  playPauseToggleEl.disabled = null;
  shuffleEl.disabled = null;
  nextEl.disabled = null;
  prevEl.disabled = null;
  window.isPlayerReady = true;
  playPauseToggleEl.innerHTML = window.playPauseIcons.play;
}

function onPlayerStateChange(event) {
  switch (event.data) {
    case YT.PlayerState.UNSTARTED:
      console.log("unstated");
      updateSongDetails();
      playPauseToggleEl.innerHTML = window.playPauseIcons.loading;
      break;
    case YT.PlayerState.ENDED:
      console.log("ended");
      playPauseToggleEl.innerHTML = window.playPauseIcons.play;
      break;
    case YT.PlayerState.PLAYING:
      console.log("playing");
      updateSongDetails();
      document.getElementById("muzikk").style.display = "block";
      document.body.style.cursor = "auto";
      playPauseToggleEl.innerHTML = window.playPauseIcons.pause;
      break;
    case YT.PlayerState.PAUSED:
      console.log("paused");
      playPauseToggleEl.innerHTML = window.playPauseIcons.play;
      break;
    case YT.PlayerState.BUFFERING:
      console.log("buffering");
      updateSongDetails();
      playPauseToggleEl.innerHTML = window.playPauseIcons.loading;
      break;
    case YT.PlayerState.CUED:
      console.log("cued");
      updateSongDetails();
      playPauseToggleEl.innerHTML = window.playPauseIcons.play;
      break;
  }
  if (event.data == YT.PlayerState.PLAYING && !timeUpdateInterval) {
    updateTimes();
  } else if (event.data !== YT.PlayerState.PLAYING && timeUpdateInterval) {
    clearInterval(timeUpdateInterval);
    timeUpdateInterval = null;
  }
}

function updateTimes() {
  timeUpdateInterval = setInterval(function () {
    if (window.player.getPlayerState() == YT.PlayerState.PLAYING) {
      if (songCurrentTimeEl) {
        songCurrentTimeEl.innerHTML = formatTime(
          window.player.getCurrentTime(),
        );
      }
      if (songSeekBarEl) {
        songSeekBarEl.value = Math.ceil(window.player.getCurrentTime());
      }
    }
  }, 1000);
}

function getTextWidth() {
  return window.innerWidth > 768 ? 35 : 15;
}

function updateSongDetails() {
  const videoData = window.player.getVideoData();
  if (videoData.title) {
    songNameEl.innerHTML = videoData.title;
    songNameEl.title = videoData.title;
    if (videoData.title.length > getTextWidth()) {
      songNameEl.parentElement.classList.add("marquee");
    } else {
      songNameEl.parentElement.classList.remove("marquee");
    }
  }
  if (videoData.author) {
    artistNameEl.innerHTML = videoData.author;
    artistNameEl.title = videoData.author;
    if (videoData.author.length > getTextWidth()) {
      artistNameEl.parentElement.classList.add("marquee");
    } else {
      artistNameEl.parentElement.classList.remove("marquee");
    }
  }
  songSeekBarEl.max = Math.ceil(window.player.getDuration());
  songSeekBarEl.value = 0;
  if (songDurationEl) {
    songDurationEl.innerHTML = formatTime(window.player.getDuration());
  }

  songImageEl.style.backgroundImage = `url("https://img.youtube.com/vi/${videoData.video_id}/0.jpg")`;
}

function stopVideo() {
  window.player.stopVideo();
  //
}

function playPauseToggle() {
  switch (window.player.getPlayerState()) {
    case -1: // unstarted
      window.player.playVideo();
      playPauseToggleEl.innerHTML = window.playPauseIcons.pause;
      break;
    case 0: // ended
      window.player.playVideo();
      playPauseToggleEl.innerHTML = window.playPauseIcons.pause;
      break;
    case 1: // playing
      window.player.pauseVideo();
      playPauseToggleEl.innerHTML = window.playPauseIcons.play;
      break;
    case 2: // paused
      window.player.playVideo();
      playPauseToggleEl.innerHTML = window.playPauseIcons.pause;
      break;
    case 3: // bufferring
      break;
    case 5: // video cued
      window.player.playVideo();
      playPauseToggleEl.innerHTML = window.playPauseIcons.pause;
      break;
  }
}

function formatTime(timeSecs) {
  timeSecs = Math.floor(timeSecs);
  const ss = Math.floor(timeSecs % 60);
  const mm = Math.floor((timeSecs / 60) % 60);
  const hh = Math.floor((timeSecs / 60 / 60) % 60);

  return `${hh > 0 ? `${formatNumber(hh)}:` : ""}${formatNumber(mm)}:${formatNumber(ss)}`;
}

function formatNumber(n) {
  return (n >= 10 ? "" : "0") + n.toString();
}

window.playYTSongById = (videoId) => {
  window.player.loadVideoById(videoId);
  document.getElementById("play").innerHTML = window.playPauseIcons.loading;
};

songSeekBarEl.addEventListener("change", (event) => {
  console.log("seeking to", Number(event.target.value));
  window.player.seekTo(Number(event.target.value));
});

document.getElementById("play").addEventListener("click", () => {
  playPauseToggle();
});
document.getElementById("next").addEventListener("click", () => {});
document.getElementById("prev").addEventListener("click", () => {});
