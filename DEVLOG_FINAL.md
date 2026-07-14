# From Legacy to Launch: How Gemini CLI Helped Me Build a Meta-Search Engine to Honor My Late Father

### Introduction: The Brink of Shutdown
In my final year of college, I built a project that was more than just code. **Kshanik Search** was a tribute to my father, Kshanik Kumar Biswas, who passed away on May 9th, 2021. I wanted to create a tool that represented his lifelong pursuit of knowledge. 

However, by late 2025, the project was on the brink of shutting down. I was relying on search engine APIs like SerpApi and Google Custom Search, but the low free-tier quotas and rising costs made it impossible to maintain as a "lasting legacy." I was stuck.

### The Turning Point: Gemini CLI and Agentic AI
This is where the **Gemini CLI** (and its upcoming successor, **Antigravity CLI**) changed everything. As the Gemini CLI prepares to sundown this June, I look back at how it became my primary collaborator.

Instead of giving up, I used Gemini to learn about **Agentic AI**—not just as a buzzword, but as a practical way to orchestrate complex tasks. Gemini didn't just write code for me; it helped me *architect* a solution from the ground up.

### Building a Meta-Search Engine from Scratch
I decided to stop relying on expensive third-party aggregators and build my own meta-search backend in Go. Using Gemini's research and execution capabilities, I implemented:
- **Concurrent Fan-out:** A backend that queries multiple providers (DuckDuckGo, Wikipedia, GitHub, Reddit, Stack Overflow) in parallel.
- **Hybrid Ranking:** A custom engine to merge, deduplicate, and rank results based on relevance and source reliability.
- **Memory Efficiency:** Optimized to run on minimal hardware.

### The "Old Laptop & ngrok" Miracle
With the 5th anniversary of my father's passing (May 9th, 2026) approaching, I set a hard deadline. I had polished the frontend and hosted it on **Netlify**, but I didn't have a production server for the backend.

In a move of sheer necessity and a bit of "hacker spirit," I ran the Go backend on my old laptop and used **ngrok** to expose it to the world. To my surprise, it worked flawlessly. The transition from a failing API-dependent site to a self-sustained meta-search engine was complete, right on time for the anniversary.

### Professional Growth: Transitioning to EY
This journey wasn't just about a project; it was about my career. Through the deep dives into system design and AI orchestration with Gemini, I transitioned from a traditional software role to a **Gen AI Engineer at EY**. 

The lessons I learned—how to guide an AI agent, how to handle state in complex LLM interactions, and how to build resilient systems—are exactly what the industry is looking for today.

### The Future: Antigravity and Beyond
As Gemini CLI transitions to **Antigravity CLI**, I'm not sad; I'm excited. This project proved that AI is the ultimate leverage. I am still learning every day, currently moving the backend to a dedicated Oracle Cloud ARM instance and planning a "GPT-like Knowledge Search" tool for human genealogy.

Kshanik Search lives on, not just as a tribute to my father, but as a testament to how far we can go when we embrace the next wave of AI.

---

*This blog was written with the help of Gemini CLI, as a final salute to the tool that helped me bridge the gap between a college dream and a professional reality.*
