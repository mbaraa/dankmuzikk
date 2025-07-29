"use strict";

const CACHE_NAME = "dank-cache-v0.2.5-offline";

let PRECACHE_ASSETS = [
  "/",
  "/about",
  "/privacy",
  "/profile",
  "/playlists",
  "/library/favorites",
  "/service-worker.js",
  "/static/android-chrome-144x144.png",
  "/static/android-chrome-192x192.png",
  "/static/android-chrome-256x256.png",
  "/static/android-chrome-36x36.png",
  "/static/android-chrome-384x384.png",
  "/static/android-chrome-48x48.png",
  "/static/android-chrome-512x512.png",
  "/static/android-chrome-72x72.png",
  "/static/android-chrome-96x96.png",
  "/static/apple-touch-icon.png",
  "/static/apple-touch-icon-precomposed.png",
  "/static/browserconfig.xml",
  "/static/css/player-seeker.css",
  "/static/css/refresher.css",
  "/static/css/tailwind.css",
  "/static/css/themes/black.css",
  "/static/css/themes/dank.css",
  "/static/css/themes/white.css",
  "/static/css/ubuntu-font.css",
  "/static/favicon-16x16.png",
  "/static/favicon-32x32.png",
  "/static/favicon.ico",
  "/static/fonts/ubuntu/Ubuntu-BoldItalic.ttf",
  "/static/fonts/ubuntu/Ubuntu-Bold.ttf",
  "/static/fonts/ubuntu/Ubuntu-Italic.ttf",
  "/static/fonts/ubuntu/Ubuntu-LightItalic.ttf",
  "/static/fonts/ubuntu/Ubuntu-Light.ttf",
  "/static/fonts/ubuntu/Ubuntu-MediumItalic.ttf",
  "/static/fonts/ubuntu/Ubuntu-Medium.ttf",
  "/static/fonts/ubuntu/Ubuntu-Regular.ttf",
  "/static/images/album-cover-icon.svg",
  "/static/images/default-pfp.svg",
  "/static/images/error-img.webp",
  "/static/images/frank-cropped.png",
  "/static/images/github.svg",
  "/static/images/google.webp",
  "/static/images/logo-big.webp",
  "/static/images/logo.webp",
  "/static/js/htmx/htmx.min.js",
  "/static/js/htmx/hyperscript.min.js",
  "/static/js/htmx/json-enc.js",
  "/static/js/htmx/loading-states.js",
  "/static/js/player_icons.js",
  "/static/js/player.js",
  "/static/js/player_shortcuts.js",
  "/static/js/player_ui.js",
  "/static/js/refresher.js",
  "/static/js/router.js",
  "/static/js/search.js",
  "/static/js/theme.js",
  "/static/js/utils.js",
  "/static/mstile-150x150.png",
  "/static/safari-pinned-tab.svg",
  "/static/site.webmanifest",
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    /**
     * @type {Cache} cache
     */
    caches.open(CACHE_NAME).then((cache) => cache.addAll(PRECACHE_ASSETS)),
  );
});

self.addEventListener("message", (event) => {
  if (event.data && event.data.type === "SKIP_WAITING") {
    self.skipWaiting();
  }
});

self.addEventListener("activate", (event) => {
  event.waitUntil(caches.keys().then((cacheNames) => Promise.all(cacheNames)));
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.open(CACHE_NAME).then((cache) =>
      cache.match(event.request).then(
        (response) =>
          response ||
          fetch(event.request).then((response) => {
            cache.put(event.request, response.clone());
            return response;
          }),
      ),
    ),
  );
});
