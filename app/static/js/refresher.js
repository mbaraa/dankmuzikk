/* I herby admit that this code is a copy-pasta from https://developer.chrome.com/blog/overscroll-behavior/ */

"use strict";

const mainContentsEl = document.getElementById("main-contents");
let _startY = 0;

async function simulateRefreshAction() {
  const sleep = (timeout) =>
    new Promise((resolve) => setTimeout(resolve, timeout));

  const transitionEnd = function (propertyName, node) {
    return new Promise((resolve) => {
      function callback(e) {
        e.stopPropagation();
        if (e.propertyName === propertyName) {
          node.removeEventListener("transitionend", callback);
          resolve(e);
        }
      }
      node.addEventListener("transitionend", callback);
    });
  };

  const refresher = document.querySelector(".refresher");

  document.body.classList.add("refreshing");
  await sleep(700);

  refresher.classList.add("shrink");
  await transitionEnd("transform", refresher);
  refresher.classList.add("done");

  refresher.classList.remove("shrink");
  document.body.classList.remove("refreshing");
  await sleep(0); // let new styles settle.
  refresher.classList.remove("done");
}

mainContentsEl.addEventListener(
  "touchstart",
  (e) => {
    _startY = e.touches[0].pageY;
  },
  { passive: true },
);

mainContentsEl.addEventListener(
  "touchmove",
  async (e) => {
    const y = e.touches[0].pageY;
    if (
      document.scrollingElement.scrollTop === 0 &&
      y > _startY &&
      !document.body.classList.contains("refreshing")
    ) {
      await simulateRefreshAction();
      await updateMainContent(window.location.pathname);
    }
  },
  { passive: true },
);
