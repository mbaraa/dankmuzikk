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
          DEFAULT: "#433C36",
          light: "#F9E6E6",
        },
        primary2: {
          DEFAULT: "#615252",
          light: "#E2BCBD",
        },
        secondary: {
          DEFAULT: "#9A7272",
          light: "#EF766B",
        },
        secondary2: {
          DEFAULT: "#C37474",
          light: "#C05B51",
        },
        accent: {
          DEFAULT: "#C05B51",
          light: "#C37474",
        },
        accent2: {
          DEFAULT: "#EF766B",
          light: "#9A7272",
        },
        accent3: {
          DEFAULT: "#E2BCBD",
          light: "#615252",
        },
        accent4: {
          DEFAULT: "#F9E6E6",
          light: "#433C36",
        },
      },
    },
  },
  plugins: [require("daisyui")],
};
