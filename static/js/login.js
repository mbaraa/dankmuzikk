"use strict";

function setContainerFullHeight(id) {
  document.getElementById(id).style.height =
    (
      window.innerHeight -
      document.getElementById("dank-header").getBoundingClientRect().height
    ).toString() + "px";
}
window.addEventListener("load", () => {
  setContainerFullHeight("frank");
  setContainerFullHeight("login-form");
});
window.addEventListener("resize", () => {
  setContainerFullHeight("frank");
  setContainerFullHeight("login-form");
});
