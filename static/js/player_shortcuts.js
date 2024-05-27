"use strict";

/**
 * Using YouTube's applicaple shortcuts: https://support.google.com/youtube/answer/7631406?hl=en
 */
const shortcuts = {
  " ": togglePP,
  k: togglePP,
  n: nextMuzikk,
  N: nextMuzikk,
  p: previousMuzikk,
  P: previousMuzikk,
  s: stopMuzikk,
  m: mute,
  M: mute,
  0: () => (audioPlayerEl.currentTime = 0),
  1: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.1),
  2: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.2),
  3: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.3),
  4: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.4),
  5: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.5),
  6: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.6),
  7: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.7),
  8: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.8),
  9: () => (audioPlayerEl.currentTime = audioPlayerEl.duration * 0.9),
  $: () => (audioPlayerEl.currentTime = audioPlayerEl.duration),
  j: () => (audioPlayerEl.currentTime -= 10),
  l: () => (audioPlayerEl.currentTime += 10),
  ArrowRight: () => setVolume(audioPlayerEl.volume + 0.1),
  ArrowLeft: () => setVolume(audioPlayerEl.volume - 0.1),
  i: expand,
  I: collapse,
  "/": () => searchInputEl.focus(),
};

/**
 * @param {KeyboardEvent} e
 */
document.addEventListener("keyup", (e) => {
  if (
    [searchFormEl, searchInputEl, searchSugEl].includes(document.activeElement)
  ) {
    return;
  }
  const action = shortcuts[e.key];
  if (action) {
    e.stopImmediatePropagation();
    e.preventDefault();
    action();
  }
});

/**
 * @param {KeyboardEvent} e
 */
document.addEventListener("keydown", (e) => {
  if (
    [searchFormEl, searchInputEl, searchSugEl].includes(document.activeElement)
  ) {
    return;
  }
  const action = shortcuts[e.key];
  if (action) {
    e.stopImmediatePropagation();
    e.preventDefault();
  }
});
