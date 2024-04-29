"use strict";

function toggleLoading() {
  if (loadingEl.classList.contains("hidden")) {
    loadingEl.classList.remove("hidden");
  } else {
    loadingEl.classList.add("hidden");
  }
}

/**
 * @param {number} timeSecs
 */
function formatTime(timeSecs) {
  timeSecs = Math.floor(timeSecs);
  const ss = Math.floor(timeSecs % 60);
  const mm = Math.floor((timeSecs / 60) % 60);
  const hh = Math.floor((timeSecs / 60 / 60) % 60);

  return `${hh > 0 ? `${formatNumber(hh)}:` : ""}${formatNumber(mm)}:${formatNumber(
    ss,
  )}`;
}

/**
 * @param {number} timeSecs
 *
 * @returns string
 */
function formatNumber(n) {
  return (n >= 10 ? "" : "0") + n.toString();
}

/**
 * @returns number
 */
function getTextWidth() {
  return window.innerWidth > 768 ? 35 : 15;
}

window.Utils = {
  toggleLoading,
  formatTime,
  formatNumber,
  getTextWidth,
};
