"use strict";

const themes = {
  dank: {
    primary: "#3E732C",
    primary20: "#3E732C33",
    primary30: "#3E732C4c",
    primary69: "#3E732Cb0",
    secondary: "#ffffff",
    secondary20: "#ffffff33",
    secondary30: "#ffffff4c",
    secondary69: "#ffffffb0",
    accent: "#000000",
    accent20: "#00000033",
    accent30: "#0000004c",
    accent69: "#000000b0",
  },
  black: {
    primary: "#000000",
    primary20: "#00000033",
    primary30: "#0000004c",
    primary69: "#000000b0",
    secondary: "#ffffff",
    secondary20: "#ffffff33",
    secondary30: "#ffffff4c",
    secondary69: "#ffffffb0",
    accent: "#236104",
    accent20: "#23610433",
    accent30: "#2361044C",
    accent69: "#236104B0",
  },
  white: {
    primary: "#ffffff",
    primary20: "#ffffff33",
    primary30: "#ffffff4c",
    primary69: "#ffffffb0",
    secondary: "#3E732C",
    secondary20: "#3E732C33",
    secondary30: "#3E732C4c",
    secondary69: "#3E732Cb0",
    accent: "#d5ffc1",
    accent20: "#d5ffc133",
    accent30: "#d5ffc14c",
    accent69: "#d5ffc1b0",
  },
};

/**
 * @param {string} themeName
 */
function changeTheme(themeName) {
  const theme = themes[themeName];
  if (!theme) {
    return;
  }
  window.Utils.setCookie("theme-name", themeName);
  const style = document.documentElement.style;

  style.setProperty("--primary-color", theme.primary);
  style.setProperty("--primary-color-20", theme.primary20);
  style.setProperty("--primary-color-30", theme.primary30);
  style.setProperty("--primary-color-69", theme.primary69);
  style.setProperty("--secondary-color", theme.secondary);
  style.setProperty("--secondary-color-20", theme.secondary20);
  style.setProperty("--secondary-color-30", theme.secondary30);
  style.setProperty("--secondary-color-69", theme.secondary69);
  style.setProperty("--accent-color", theme.accent);
  style.setProperty("--accent-color-20", theme.accent20);
  style.setProperty("--accent-color-30", theme.accent30);
  style.setProperty("--accent-color-69", theme.accent69);
  //document.getElementById("popover-theme-switcher").style.display = "none";
}

(() => {
  const userTheme = window.Utils.getCookie("theme-name");
  if (userTheme) {
    changeTheme(userTheme);
    return;
  }

  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: light)").matches
  ) {
    changeTheme("white");
  }

  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    changeTheme("black");
  }
})();

window
  .matchMedia("(prefers-color-scheme: dark)")
  .addEventListener("change", (e) => {
    changeTheme(e.matches ? "black" : "white");
  });

window.Theme = { changeTheme };
