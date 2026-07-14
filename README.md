# 🌟 Kshanik Search

A heartfelt tribute to **Kshanik Kumar Biswas (1943–2021)**, this search engine blends modern technology with a meaningful legacy. Explore the web with purpose and discover the story behind this project.

---

## 📖 About

Kshanik Search is more than a search engine—it's a dedication to the enduring memory of Kshanik Kumar Biswas. With a sleek interface and thoughtful features, this project combines functionality with a mission to honor a remarkable life.

---

## ✨ Features

- **Powerful Search**: Effortlessly explore the web with our custom-built search engine.
- **Dark Mode Toggle**: Switch between light and dark themes for a comfortable experience.
- **About Page**: Dive into the heartfelt inspiration behind Kshanik Search.

---

## 🚀 Upcoming Features (Roadmap)

- [ ] **GPT-Powered Search**: Supercharge results with AI-driven GPT responses.
- [ ] **AI Voice Assistant**: Enable hands-free searches with an intelligent voice interface.

---

## 🛠️ How to Run

Get started with Kshanik Search in just a few steps:

1. Clone the repository.
2. Install dependencies: `npm install`
3. Start the development server: `npm run dev`
4. Build for production: `npm run build`

## Deployment
This project is deployed on Netlify. Visit [Kshanik Search](https://kshaniksearch.netlify.app) to explore.


## Netlify + Go backend integration (fix for blank results)

This frontend is static and cannot run the Go API inside Netlify directly. To show results correctly:

1. Deploy the Go backend (`cmd/server`) to a backend host (Render/Fly.io/Railway/etc).
2. Set a Netlify redirect so `/api/*` points to your Go service.
3. Set frontend env var `VITE_SEARCH_API_BASE_URL`:
   - `/api` (recommended with Netlify redirect), or
   - direct backend URL (for local debugging).
4. Ensure CORS is enabled on the Go API if you use a direct URL instead of `/api`.

The app now fetches from `GET /search?q=` and expects normalized JSON array objects:
`title`, `snippet`, `url`, `source`, `score`, `timestamp`.
