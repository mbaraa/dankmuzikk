"use strict";

/**
 * Using YouTube's applicable shortcuts: https://support.google.com/youtube/answer/7631406?hl=en
 *
 * set those characters in Vimium's excluded patterns if you still wanna navigate the site using Vimium,
 * l r k n N p P s m M 0 1 2 3 4 5 6 7 8 9 $ i I /
 */
const shortcuts = {
  " ": togglePP,
  l: () => toggleLoop(),
  r: () => toggleShuffle(),
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
  ArrowLeft: () => (audioPlayerEl.currentTime -= 10),
  ArrowRight: () => (audioPlayerEl.currentTime += 10),
  ArrowUp: () => setVolume(audioPlayerEl.volume + 0.1),
  ArrowDown: () => setVolume(audioPlayerEl.volume - 0.1),
  i: expand,
  I: collapse,
  "/": () => searchInputEl.focus(),
};

/**
 * @param {KeyboardEvent} e
 */
document.addEventListener("keyup", (e) => {
  if (
    document.activeElement.tagName === "INPUT" ||
    document.activeElement.id.startsWith("search-suggestion")
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
  if (document.activeElement.tagName === "INPUT") {
    return;
  }
  const action = shortcuts[e.key];
  if (action) {
    e.stopImmediatePropagation();
    e.preventDefault();
  }
});
