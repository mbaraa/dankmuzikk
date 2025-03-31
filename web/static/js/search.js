"use strict";

const searchInputEl = document.getElementById("search-input");

/**
 * @param {KeyboardEvent} e
 */
function moveInSuggestions(e) {
  const t = e.target;
  if (!t) return;
  const p = t.parentElement;
  if (!p) return;
  const next = p.nextSibling;
  const prev = p.previousSibling;

  if (e.key === "ArrowDown") {
    e.preventDefault();
    if (next && next.firstChild) {
      next.firstChild.focus();
      searchInputEl.value = next.firstChild.innerText;
    }
  }
  if (e.key === "ArrowUp") {
    e.preventDefault();
    if (prev && prev.firstChild) {
      searchInputEl.value = prev.firstChild.innerText;
      prev.firstChild.focus();
    } else {
      searchInputEl.focus();
    }
  }
}

window.Search = {
  moveInSuggestions,
};
