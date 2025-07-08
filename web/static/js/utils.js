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
 * @param {number} timeMs
 */
function formatTimeMs(timeMs) {
  timeMs = Math.floor(timeMs * 1000);
  const ms = Math.floor(timeMs % 100);
  timeMs = Math.floor(timeMs / 1000);
  const ss = Math.floor(timeMs % 60);
  const mm = Math.floor((timeMs / 60) % 60);
  const hh = Math.floor((timeMs / 60 / 60) % 60);

  return `${hh > 0 ? `${formatNumber(hh)}:` : ""}${formatNumber(mm)}:${formatNumber(ss)}.${formatNumber(ms)}`;
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

/**
 * @param {string} key
 */
function getCookie(key) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${key}=`);
  if (parts.length === 2) return parts.pop().split(";").shift();
}

/**
 * @param {string} key
 * @param {string} value
 */
function setCookie(key, value) {
  const date = new Date();
  date.setTime(date.getTime() + 365 * 24 * 60 * 60 * 1000);
  const expires = "; expires=" + date.toUTCString();
  document.cookie = key + "=" + value + expires + "; path=/";
}

async function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

window.Utils = {
  showLoading,
  hideLoading,
  formatTime,
  formatTimeMs,
  formatNumber,
  getTextWidth,
  copyTextToClipboard,
  getCookie,
  setCookie,
  sleep,
};
