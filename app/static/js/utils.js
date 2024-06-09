"use strict";

const loadingEl = document.getElementById("loading");

function showLoading() {
  loadingEl.classList.remove("hidden");
}

function hideLoading() {
  loadingEl.classList.add("hidden");
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

/**
 * @param {string} text
 */
function copyTextToClipboard(text) {
  const textArea = document.getElementById("clipboard-text-blyat");
  textArea.hidden = false;
  textArea.value = text;
  textArea.select();
  textArea.setSelectionRange(0, 300);
  document.execCommand("copy");
  textArea.hidden = true;
}

let isMobile = window.innerWidth < 768;

/**
 * @param {EventTarget<Window>} e
 */
window.addEventListener("resize", (e) => {
  if (e.target.innerWidth < 768 && !isMobile) {
    isMobile = true;
    window.location.reload();
  }
  if (e.target.innerWidth > 768 && isMobile) {
    isMobile = false;
    window.location.reload();
  }
});

window.Utils = {
  showLoading,
  hideLoading,
  formatTime,
  formatNumber,
  getTextWidth,
  copyTextToClipboard,
};
