"use strict";

const playerButtonsIcons = {
  playDisabled: ` <img class="w-[50px] h-[50px]" src="/static/images/play-disabled-icon.svg" alt="Play"/>`,
  play: `<img class="w-[50px] h-[50px]" src="/static/images/play-icon.svg" alt="Play"/>`,
  pauseDisabled: `<img class="w-[50px] h-[50px]" src="/static/images/pause-disabled-icon.svg" alt="Pause"/>`,
  pause: `<img class="w-[50px] h-[50px]" src="/static/images/pause-icon.svg" alt="Pause"/>`,
  loop: `<img class="w-[40px]" src="/static/images/loop-icon.svg" alt="Loop"/>`,
  loopOnce: `<img class="w-[40px]" src="/static/images/loop-once-icon.svg" alt="Loop Once"/>`,
  loopOff: `<img class="w-[40px]" src="/static/images/loop-off-icon.svg" alt="Loop Off"/>`,
  loading: `<div class="loader !h-10 !w-10"></div>`,
};

const loopModes = [
  { icon: "loop-off-icon.svg", mode: "OFF" },
  { icon: "loop-once-icon.svg", mode: "ONCE" },
  // TODO: implement this
  //{ icon: "loop-icon.svg", mode: "ALL"},
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
  loadingEl = document.getElementById("loading"),
  audioPlayerEl = document.getElementById("audio-player"),
  muzikkContainerEl = document.getElementById("muzikk");

let currentLoopIdx = 0;

/**
 * @param {{id: string, artist: string, thumbnailUrl: string, title: string}} videoData
 */
function setMediaSession(videoData) {
  if (!("mediaSession" in navigator)) {
    console.error("Browser doesn't support mediaSession");
    return;
  }
  navigator.mediaSession.metadata = new MediaMetadata({
    title: videoData.title,
    artist: videoData.artist,
    album: videoData.artist,
    artwork: [
      "96x96",
      "128x128",
      "192x192",
      "256x256",
      "384x384",
      "512x512",
    ].map((i) => {
      return {
        src: videoData.thumbnailUrl,
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

async function fetchMusic(youtubeId) {
  playPauseToggleEl.innerHTML = playerButtonsIcons.loading;
  document.body.style.cursor = "progress";
  Utils.showLoading();

  await fetch("/api/song/download/" + youtubeId)
    .then((res) => console.log(res))
    .catch((err) => console.error(err));

  if (audioPlayerEl) {
    stopMuzikk();
  }
  audioPlayerEl.src = `/music/${youtubeId}.mp3`;
  audioPlayerEl.load();
}

async function playYTSongById(id, thumbnailUrl, title, artist) {
  const videoData = { id, thumbnailUrl, title, artist };
  await fetchMusic(videoData.id);
  setMediaSession(videoData);
  showPlayer();

  if (videoData.title) {
    songNameEl.innerHTML = videoData.title;
    songNameEl.title = videoData.title;
    if (videoData.title.length > Utils.getTextWidth()) {
      songNameEl.parentElement.classList.add("marquee");
    } else {
      songNameEl.parentElement.classList.remove("marquee");
    }
  }
  if (videoData.artist) {
    artistNameEl.innerHTML = videoData.artist;
    artistNameEl.title = videoData.artist;
    if (videoData.artist.length > Utils.getTextWidth()) {
      artistNameEl.parentElement.classList.add("marquee");
    } else {
      artistNameEl.parentElement.classList.remove("marquee");
    }
  }

  playMuzikk();
  songImageEl.style.backgroundImage = `url("${videoData.thumbnailUrl}")`;
}

function showPlayer() {
  muzikkContainerEl.style.display = "block";
}

function hidePlayer() {
  muzikkContainerEl.style.display = "none";
  audioPlayerEl.stopMuzikk();
}

loopEl.addEventListener("click", (event) => {
  currentLoopIdx = (currentLoopIdx + 1) % loopModes.length;
  event.target.src = "/static/images/" + loopModes[currentLoopIdx].icon;
});

playPauseToggleEl.addEventListener("click", () => {
  playPauseToggle();
});

nextEl.addEventListener("click", () => {});

prevEl.addEventListener("click", () => {});

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
});

audioPlayerEl.addEventListener("ended", () => {
  switch (loopModes[currentLoopIdx].mode) {
    case "OFF":
      stopMuzikk();
      break;
    case "ONCE":
      stopMuzikk();
      playMuzikk();
      break;
    case "ALL":
      break;
  }
});

audioPlayerEl.addEventListener("progress", () => {
  console.log("downloading...");
});

window.Player = {
  playYTSongById,
  showPlayer,
  hidePlayer,
};
