/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        Ubuntu: ["Ubuntu", "sans"],
      },
      colors: {
        primary: {
          DEFAULT: "var(--primary-color)",
          trans: {
            20: "var(--primary-color-20)",
          },
        },
        secondary: {
          DEFAULT: "var(--secondary-color)",
          trans: {
            20: "var(--secondary-color-20)",
          },
        },
        accent: {
          DEFAULT: "var(--accent-color)",
          trans: {
            20: "var(--accent-color-20)",
            30: "var(--accent-color-30)",
            69: "var(--accent-color-69)",
          },
        },
      },
    },
  },
  plugins: [],
};
