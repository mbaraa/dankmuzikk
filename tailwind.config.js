/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/**/*.templ"],
  theme: {
    extend: {
      fontFamily: {
        AudioNugget: ["Audio Nugget"],
        MyriadPro: ["Myriad Pro"],
      },
      colors: {
        primary: {
          DEFAULT: "#7ACB65",
        },
        secondary: {
          DEFAULT: "#9EE07E",
        },
        accent: {
          DEFAULT: "#4C8C36",
        },
        black: {
          DEFAULT: "#000",
          trans: {
            100: "#00000044",
            200: "#00000055",
            400: "#00000088",
          },
        },
      },
    },
  },
  plugins: [],
};
