# 🚀 Oracle Cloud Deployment Guide (Advanced)

This guide provides steps for deploying the redesigned Hybrid Search Engine on **Oracle Cloud Infrastructure (OCI)** using modern containerization and CI/CD.

---

## Phase 1: Infrastructure Setup (OCI ARM)

1.  **Instance:** Create a **VM.Standard.A1.Flex** instance (Ubuntu 24.04).
    *   Recommended: **4 OCPU / 24GB RAM**.
2.  **VCN Security List:**
    *   Open **80 (HTTP)**, **443 (HTTPS)**.
    *   Open **6379 (Redis)** - *Only for internal VCN traffic*.
    *   Open **7700 (Meilisearch)** - *Only for internal VCN traffic*.

---

## Phase 2: Server Preparation

SSH into your server and run the following to install the modern stack:

```bash
# Install Docker and Docker Compose
sudo apt update
sudo apt install -y docker.io docker-compose
sudo usermod -aG docker ubuntu
# (Log out and back in for group changes)

# Install Caddy for SSL
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy
```

---

## Phase 3: Deployment via Docker Compose

Create a `docker-compose.yml` in `/home/ubuntu/kshanik_search/` to manage the hybrid stack.

```yaml
version: '3.8'
services:
  backend:
    build: 
      context: ./backend
      dockerfile: Dockerfile
    restart: always
    env_file: .env
    ports:
      - "8080:8080"
    depends_on:
      - redis
      - meilisearch

  redis:
    image: redis:7-alpine
    restart: always

  meilisearch:
    image: getmeili/meilisearch:latest
    restart: always
    volumes:
      - ./meili_data:/meili_data
    environment:
      - MEILI_MASTER_KEY=your_secure_key_here
```

---

## Phase 4: CI/CD with GitHub Actions

To automate deployment, add a GitHub Action `.github/workflows/deploy.yml`:

1.  **Build Docker Image:** Build natively for `linux/arm64`.
2.  **Push to OCI Registry (OCIR):** Or use Docker Hub.
3.  **Deploy via SSH:**
    ```yaml
    - name: Deploy to OCI
      uses: appleboy/ssh-action@master
      with:
        host: ${{ secrets.OCI_IP }}
        username: ubuntu
        key: ${{ secrets.OCI_SSH_KEY }}
        script: |
          cd ~/kshanik_search
          git pull
          docker-compose up -d --build
    ```

---

## Phase 5: Reverse Proxy (Caddy)

Configure `Caddyfile` for automatic SSL:

```caddyfile
search.yourdomain.com {
    reverse_proxy localhost:8080
}
```

---

## ⚡ Computed Advantages of the Redesign

1.  **Atomic Deployment:** Docker Compose ensures all services (Redis, Index, Backend) start in the correct order.
2.  **Zero-Downtime:** Using Docker's health checks and Caddy's hot-reloading ensures the search engine stays "Alive".
3.  **Scalable Crawling:** The Redis service allows adding more ARM instances as "Crawler Workers" without changing the main backend.
4.  **OCI Cost Efficiency:** Utilizing the full 24GB of RAM for the local index allows for millions of documents to be indexed for free.


│  To resume this session: gemini --resume 'f8f50734-3e3b-4829-a851-d36265f3de46'                    │