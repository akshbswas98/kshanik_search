/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: "#003527",
        "primary-container": "#064e3b",
        "on-primary": "#ffffff",
        secondary: "#515f74",
        "secondary-container": "#d5e3fd",
        surface: "#f7f9fb",
        "on-surface": "#191c1e",
        "on-surface-variant": "#404944",
        emerald: {
          700: '#065f46',
          800: '#064e3b',
          900: '#022c22',
        },
        slate: {
          50: '#f8fafc',
          100: '#f1f5f9',
          200: '#e2e8f0',
          600: '#475569',
          800: '#1e293b',
          900: '#0f172a',
        }
      },
      fontFamily: {
        serif: ['"Noto Serif"', 'Georgia', 'serif'],
        sans: ['Manrope', 'Inter', 'sans-serif'],
      }
    },
  },
  plugins: [],
}
