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
        },
        secondary: {
          DEFAULT: "var(--secondary-color)",
        },
        accent: {
          DEFAULT: "var(--accent-color)",
        },
        black: {
          DEFAULT: "#000",
          trans: {
            100: "#00000044",
            200: "#00000055",
            300: "#00000066",
            400: "#00000088",
            500: "#000000B0",
          },
        },
        white: {
          DEFAULT: "#fff",
          trans: {
            100: "#FFFFFF25",
            200: "#FFFFFFEE",
          },
        },
      },
    },
  },
  plugins: [],
};
