"use strict";

const searchFormEl = document.getElementById("search-form"),
  searchInputEl = document.getElementById("search-input"),
  searchSugEl = document.getElementById("search-suggestions-container");

let focusedSuggestionIndex = 0;

function searchNoReload(searchQuery) {
  searchSugEl.innerText = "";
  searchFormEl.blur();
  searchInputEl.blur();
  searchInputEl.value = searchQuery;
}

searchFormEl.addEventListener("submit", (e) => {
  e.preventDefault();
  searchNoReload(e.target.query.value);
});

document.getElementById("search-icon").addEventListener("click", () => {
  searchNoReload(searchFormEl.query.value);
});

searchInputEl.addEventListener("keydown", (e) => {
  if (e.key !== "ArrowDown") {
    return;
  }
  moveToSuggestions();
});

function moveToSuggestions() {
  let searchSuggestionsEl = document.getElementById(
    "search-suggestion-" + focusedSuggestionIndex.toString(),
  );
  // sometimes it needs a second to initialize.
  if (!searchSuggestionsEl) {
    searchSuggestionsEl = document.getElementById(
      "search-suggestion-" + focusedSuggestionIndex.toString(),
    );
  }
  if (!searchSuggestionsEl) {
    return;
  }
  const moveToNextSuggestion = (e) => {
    const numSuggestions = (
      document.getElementById("search-suggestions").children ?? []
    ).length;
    if (e.key === "ArrowDown") {
      focusedSuggestionIndex = (focusedSuggestionIndex + 1) % numSuggestions;
      moveToSuggestions();
      searchSuggestionsEl.removeEventListener("keydown", moveToNextSuggestion);
    }
    if (e.key === "ArrowUp") {
      focusedSuggestionIndex--;
      if (focusedSuggestionIndex < 0) {
        focusedSuggestionIndex = 0;
        searchInputEl.focus();
        return;
      }
      moveToSuggestions();
      searchSuggestionsEl.removeEventListener("keydown", moveToNextSuggestion);
    }
  };
  searchSuggestionsEl.focus();
  searchSuggestionsEl.addEventListener("keydown", moveToNextSuggestion);
}

window.Search = {
  searchNoReload,
};
