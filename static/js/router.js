"use strict";

const links = [
  { check: (l) => l === "/", element: document.getElementById("/") },
  { check: (l) => l === "/about", element: document.getElementById("/about") },
  {
    check: (l) => l === "/profile",
    element: document.getElementById("/profile"),
  },
  {
    check: (l) => l.startsWith("/playlist"),
    element: document.getElementById("/playlists"),
  },
];

function updateActiveNavLink() {
  for (const link of links) {
    if (link.check(window.location.pathname)) {
      link.element.classList.add("bg-accent-trans-20");
    } else {
      link.element.classList.remove("bg-accent-trans-20");
    }
  }
}

window.addEventListener("load", () => {
  updateActiveNavLink();
});

document.addEventListener("htmx:afterRequest", function (e) {
  console.log("lol", e);
});

window.Router = { updateActiveNavLink };
