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

function menuer() {
  const menus = [];

  /**
   * @param {string} id
   */
  const __registerPopover = (id) => {
    menus.push(id);
    document.body.addEventListener("click", __remove);
  };

  /**
   * @param {MouseEvent} e
   */
  const __remove = (e) => {
    for (let i = 0; i < menus.length; i++) {
      const el = document.getElementById(menus[i]);
      if (!el) {
        continue;
      }
      const rect = el.getBoundingClientRect();
      const parentRect = el.parentElement.getBoundingClientRect();
      if (
        e.clientX < rect.left ||
        e.clientX > rect.right ||
        e.clientY + parentRect.height + 5 < rect.top ||
        e.clientY > rect.bottom
      ) {
        menus.splice(i, 1);
        el.style.display = "none";
      }
    }
  };

  return [__registerPopover];
}

const [registerPopover, registerMobileMenu, registerPopup] = menuer();

window.Utils = {
  showLoading,
  hideLoading,
  formatTime,
  formatNumber,
  getTextWidth,
  copyTextToClipboard,
  registerPopover,
};
