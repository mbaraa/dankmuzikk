"use strict";

const searchFormEl = document.getElementById("search-form"),
  searchInputEl = document.getElementById("search-input");

let focusedSuggestionIndex = 0;

searchFormEl.addEventListener("submit", (e) => {
  e.preventDefault();
  const query = encodeURIComponent(e.target.query.value);
  window.open(`/search?query=${query}`, "_self");
});

document.getElementById("search-icon").addEventListener("click", () => {
  const query = encodeURIComponent(searchFormEl.query.value);
  window.open(`/search?query=${query}`, "_self");
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
