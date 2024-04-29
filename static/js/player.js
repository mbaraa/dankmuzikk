"use strict";

const playPauseIcons = {
  playDisabled: ` <img class="w-[50px] h-[50px]" src="/static/images/play-disabled-icon.svg" alt="Play"/>`,
  play: `<img class="w-[50px] h-[50px]" src="/static/images/play-icon.svg" alt="Play"/>`,
  pauseDisabled: `<img class="w-[50px] h-[50px]" src="/static/images/pause-disabled-icon.svg" alt="Pause"/>`,
  pause: `<img class="w-[50px] h-[50px]" src="/static/images/pause-icon.svg" alt="Pause"/>`,
  loop: `<img class="w-[40px]" src="/static/images/loop-icon.svg" alt="Loop"/>`,
  loopOnce: `<img class="w-[40px]" src="/static/images/loop-once-icon.svg" alt="Loop Once"/>`,
  loopOff: `<img class="w-[40px]" src="/static/images/loop-off-icon.svg" alt="Loop Off"/>`,
  loading: `<div class="loader !h-10 !w-10"></div>`,
};

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
  audioPlayerEl = document.getElementById("audio-player");

let currentTime = 0,
  duration = 0,
  currentAudio = {},
  paused = false;

function setMediaSession(videoData) {
  if ("mediaSession" in navigator) {
    navigator.mediaSession.metadata = new MediaMetadata({
      title: videoData.title,
      artist: videoData.artist,
      album: videoData.artist,
      artwork: [
        {
          src: videoData.thumbnailUrl,
          sizes: "96x96",
          type: "image/png",
        },
      ],
    });

    navigator.mediaSession.setActionHandler("play", () => {
      playPauseToggle();
    });
    navigator.mediaSession.setActionHandler("pause", () => {
      playPauseToggle();
    });
    navigator.mediaSession.setActionHandler("stop", () => {
      audioPlayerEl.pause();
      audioPlayerEl.currentTime = 0;
      document.getElementById("play").innerHTML = playPauseIcons.play;
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
      if (audioPlayerEl.currentTime + seekTo > duration) {
        seekTo = 0;
      }
      audioPlayerEl.currentTime += seekTo;
    });
    navigator.mediaSession.setActionHandler("seekto", (a) => {
      const seekTime = Number(a.seekOffset);
      audioPlayerEl.currentTime = seekTime;
    });
    navigator.mediaSession.setActionHandler("previoustrack", () => {
      // previous();
    });
    navigator.mediaSession.setActionHandler("nexttrack", () => {
      // next();
    });
  }
}

window.isPlayerReady = false;

function playPauseToggle() {
  if (audioPlayerEl.paused) {
    audioPlayerEl.play();
    document.getElementById("play").innerHTML = playPauseIcons.pause;
  } else {
    audioPlayerEl.pause();
    document.getElementById("play").innerHTML = playPauseIcons.play;
  }
}

let timeUpdateInterval;

async function fetchMusic(youtubeId) {
  document.getElementById("play").innerHTML = playPauseIcons.loading;
  document.body.style.cursor = "progress";

  await fetch("/api/song/download/" + youtubeId)
    .then((res) => console.log(res))
    .catch((err) => console.error(err));

  if (audioPlayerEl) {
    audioPlayerEl.pause();
    audioPlayerEl.currentTime = 0;
    currentTime = 0;
  }
  document.getElementById("muzikk").style.display = "block";
  audioPlayerEl.src = `/music/${youtubeId}.mp3`;
  paused = false;
  audioPlayerEl.load();
  //	    	currentAudio = music;
  //	    	pageTitle = currentAudio.title +
  //translate(TranslationKeys.TITLE_PLAYING_SUFFIX);
}

function formatTime(timeSecs) {
  timeSecs = Math.floor(timeSecs);
  const ss = Math.floor(timeSecs % 60);
  const mm = Math.floor((timeSecs / 60) % 60);
  const hh = Math.floor((timeSecs / 60 / 60) % 60);

  return `${hh > 0 ? `${formatNumber(hh)}:` : ""}${formatNumber(mm)}:${formatNumber(
    ss,
  )}`;
}

function formatNumber(n) {
  return (n >= 10 ? "" : "0") + n.toString();
}

function getTextWidth() {
  return window.innerWidth > 768 ? 35 : 15;
}

window.playYTSongById = async (id, thumbnailUrl, title, artist) => {
  const videoData = { id, thumbnailUrl, title, artist };
  // window.player.loadVideoById(videoId);
  await fetchMusic(videoData.id);
  setMediaSession(videoData);
  if (videoData.title) {
    songNameEl.innerHTML = videoData.title;
    songNameEl.title = videoData.title;
    if (videoData.title.length > getTextWidth()) {
      songNameEl.parentElement.classList.add("marquee");
    } else {
      songNameEl.parentElement.classList.remove("marquee");
    }
  }
  if (videoData.artist) {
    artistNameEl.innerHTML = videoData.artist;
    artistNameEl.title = videoData.artist;
    if (videoData.artist.length > getTextWidth()) {
      artistNameEl.parentElement.classList.add("marquee");
    } else {
      artistNameEl.parentElement.classList.remove("marquee");
    }
  }

  audioPlayerEl.play();

  songImageEl.style.backgroundImage = `url("${videoData.thumbnailUrl}")`;
};

document.getElementById("play").addEventListener("click", () => {
  playPauseToggle();
});
document.getElementById("next").addEventListener("click", () => {});
document.getElementById("prev").addEventListener("click", () => {});

const loopModes = [
  { icon: "loop-off-icon.svg", mode: "OFF" },
  { icon: "loop-once-icon.svg", mode: "ONCE" },
  // TODO: implement this
  //{ icon: "loop-icon.svg", mode: "ALL"},
];
let currentLoopIdx = 0;
loopEl.addEventListener("click", (event) => {
  currentLoopIdx = (currentLoopIdx + 1) % loopModes.length;
  event.target.src = "/static/images/" + loopModes[currentLoopIdx].icon;
});

songSeekBarEl.addEventListener("change", (event) => {
  const seekTime = Number(event.target.value);
  // const audio = document.getElementById("aud") as any;
  audioPlayerEl.currentTime = seekTime;
  // player = audio;
});

audioPlayerEl.addEventListener("loadeddata", (event) => {
  playPauseToggleEl.disabled = null;
  shuffleEl.disabled = null;
  nextEl.disabled = null;
  prevEl.disabled = null;
  loopEl.disabled = null;
  window.isPlayerReady = true;

  const duration = Math.ceil(
    Number.isNaN(event.target.duration) ? 0 : event.target.duration,
  );
  songSeekBarEl.max = Math.ceil(duration);
  songSeekBarEl.value = 0;
  if (songDurationEl) {
    songDurationEl.innerHTML = formatTime(duration);
  }

  document.getElementById("play").innerHTML = playPauseIcons.pause;
  document.body.style.cursor = "auto";
  // duration = a.duration;
});

audioPlayerEl.addEventListener("timeupdate", (event) => {
  const currentTime = Math.floor(event.target.currentTime);
  if (songCurrentTimeEl) {
    songCurrentTimeEl.innerHTML = formatTime(currentTime);
  }
  if (songSeekBarEl) {
    songSeekBarEl.value = Math.ceil(currentTime);
  }
});

audioPlayerEl.addEventListener("ended", (event) => {
  switch (loopModes[currentLoopIdx].mode) {
    case "OFF":
      document.getElementById("play").innerHTML = playPauseIcons.play;
      audioPlayerEl.currentTime = 0;
      break;
    case "ONCE":
      document.getElementById("play").innerHTML = playPauseIcons.pause;
      audioPlayerEl.currentTime = 0;
      audioPlayerEl.play();
      break;
    case "ALL":
      break;
  }
});

audioPlayerEl.addEventListener("progress", (event) => {
  console.log("downloading...");
});
