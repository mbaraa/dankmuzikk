"use strict";

const mainContentsEl = document.getElementById("main-contents");

const links = [
  {
    check: (l) => l === "/",
    elements: [
      document.getElementById("/"),
      document.getElementById("/?mobile"),
    ],
  },
  {
    check: (l) => l === "/about",
    elements: [
      document.getElementById("/about"),
      document.getElementById("/about?mobile"),
    ],
  },
  {
    check: (l) => l === "/profile",
    elements: [
      document.getElementById("/profile"),
      document.getElementById("/profile?mobile"),
    ],
  },
  {
    check: (l) => l.startsWith("/playlist"),
    elements: [
      document.getElementById("/playlists"),
      document.getElementById("/playlists?mobile"),
    ],
  },
];

function updateActiveNavLink() {
  for (const link of links) {
    if (link.check(window.location.pathname)) {
      link.elements.forEach((e) => e?.classList.add("bg-accent-trans-20"));
    } else {
      link.elements.forEach((e) => e?.classList.remove("bg-accent-trans-20"));
    }
  }
}

window.addEventListener("load", () => {
  updateActiveNavLink();
});

/**
 * @param {string} path the requested path to update.
 */
async function updateMainContent(path) {
  Utils.showLoading();
  const query = new URLSearchParams(location.search);
  query.set("no_layout", "true");
  htmx
    .ajax("GET", path + "?" + query.toString(), {
      target: "#main-contents",
      swap: "innerHTML",
    })
    .catch(() => {
      window.location.reload();
    })
    .finally(() => {
      Utils.hideLoading();
      updateActiveNavLink();
    });
}

window.addEventListener("popstate", async (e) => {
  const mainContentsEl = document.getElementById("main-contents");
  if (!!mainContentsEl && !!e.target.location.pathname) {
    e.stopImmediatePropagation();
    e.preventDefault();
    await updateMainContent(e.target.location.pathname);
    return;
  }
});

document.addEventListener("htmx:afterRequest", function (e) {
  if (!!e.detail && !!e.detail.xhr) {
    const newTitle = e.detail.xhr.getResponseHeader("HX-Title");
    if (newTitle) {
      document.title = newTitle + " - DankMuzikk";
    }
  }
});

window.Router = { updateActiveNavLink, updateMainContent };
