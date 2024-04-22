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
      },
    },
  },
  plugins: [],
};
