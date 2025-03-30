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

function menuer() {
  /**
   * @type {HTMLDivElement}
   */
  let lastEl = null;

  /**
   * @param {string} id
   */
  const __registerPopover = (id) => {
    if (!id) {
      return;
    }
    if (!!lastEl) {
      lastEl.classList.add("hidden");
      document.body.removeEventListener("click", __removePopover);
    }
    lastEl = document.getElementById(id);
    if (!lastEl) {
      return;
    }
    document.body.addEventListener("click", __removePopover);
  };

  /**
   * @param {MouseEvent} e
   */
  const __removePopover = (e) => {
    const rect = lastEl.getBoundingClientRect();
    const parentRect = lastEl.parentElement.getBoundingClientRect();
    let popupChild = null;
    for (const c of lastEl.children.item(0)?.children) {
      if (c.tagName === "DIALOG") {
        popupChild = c;
      }
    }
    if (document.activeElement === popupChild) {
      return;
    }
    if (
      e.clientX < rect.left ||
      e.clientX > rect.right ||
      e.clientY + parentRect.height + 5 < rect.top ||
      e.clientY > rect.bottom + parentRect.height + 5
    ) {
      lastEl.classList.add("hidden");
      lastEl = null;
      document.body.removeEventListener("click", __removePopover);
    }
  };

  return [__registerPopover];
}

const [registerPopover, registerMobileMenu, registerPopup] = menuer();

/**
 * @param {() => Promise<any>} func
 * @param {number} times
 *
 * @returns Promise<any>
 */
async function retryer(func, times = 3) {
  try {
    return await func();
  } catch (err) {
    if (times > 0) {
      console.log("retrying ", times);
      await sleep(3500);
      return await retryer(func, times - 1);
    }
    return err;
  }
}

async function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

window.Utils = {
  showLoading,
  hideLoading,
  formatTime,
  formatNumber,
  getTextWidth,
  copyTextToClipboard,
  registerPopover,
  getCookie,
  setCookie,
  retryer,
  sleep,
};
