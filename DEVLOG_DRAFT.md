# Dev.to Post Draft: Kshanik Search Journey

Title: From Legacy to Launch: How I Rebuilt a Meta-Search Engine to Honor My Late Father
Tags: #go #react #systemdesign #ai #webdev

---

**The Genesis (2021)**
In December 2021, I built the first version of Kshanik Search. It was a simple project, but it carried a heavy purpose: a tribute to my late father, Kshanik Kumar Biswas. The goal was to create a tool that represented his pursuit of knowledge.

**The Technical Wall**
Originally, I relied on Serp API for results. While it worked, the low free-tier limits and cost became a bottleneck for a project intended to be a lasting legacy. I needed more control, more scale, and better performance.

**The Revamp: Enter Gemini CLI & Antigravity**
Earlier this year, I decided to overhaul everything. Using the Gemini CLI and Antigravity, I fundamentally redesigned the UI/UX to feel modern and "alive." But the real work was under the hood. I moved away from single-provider APIs and applied System Design principles to build a custom Meta-Search Engine in Go.

**The Architecture:**
- Concurrent Fan-out: The Go backend queries multiple providers (DuckDuckGo, Wikipedia, GitHub, etc.) in parallel.
- Hybrid Ranking: I implemented a custom ranking engine to merge and deduplicate results.
- Resource Constraints: Due to time and cost, I had to be creative.

**The Launch (May 9th)**
May 9th marks my father's death anniversary. I set a hard deadline to launch the revamp on this day. With minutes on the clock and no production server ready, I ran the backend on ngrok using my old laptop.

To my surprise, it ran like magic. The community feedback was overwhelming, proving that performance isn't always about expensive clusters—it’s about efficient code and heart.

**What’s Next? (Part 2)**
The "Old Laptop & ngrok" phase was a success, but now I’m moving to Stage 2: Deploying to a dedicated platform (Oracle Cloud ARM) to support a GPT-like Knowledge Search tool for Human Genealogy.
