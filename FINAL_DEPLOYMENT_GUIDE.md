# Final Deployment Guide: Kshanik Search Engine

This document provides end-to-end instructions for deploying the **Go Meta-Search Backend** to Oracle Cloud and the **React Frontend** to Netlify or Vercel.

---

## Part 1: Pre-Deployment Verification

Before deploying, ensure your frontend is configured to communicate with the production backend.

### 1. Update Frontend API Endpoint Handling

Since your local environment uses a Vite proxy (`/api` -> `localhost:8080`), we need to verify your React app knows how to call the production server once deployed. 

In `src/contexts/ResultsContextProvider.jsx`, your code dynamically decides the endpoint:

```javascript
const API_BASE = import.meta.env.VITE_BACKEND_URL || '';
// In development, VITE_BACKEND_URL is empty, so it uses /api/search (caught by Vite Proxy).
// In production, VITE_BACKEND_URL is set, so it uses https://api.yourdomain.com/search.
const endpoint = API_BASE ? `${API_BASE}/search?${params.toString()}` : `/api/search?${params.toString()}`;
const response = await fetch(endpoint);
```

*(I have proactively updated your `ResultsContextProvider.jsx` to support this behavior!)*

---

## Part 2: Backend Deployment (Oracle Cloud ARM)

Oracle provides a generous free tier of up to 4 ARM OCPUs and 24GB RAM, ideal for running our lightweight Go server.

### 1. Provision the Instance
1. Go to **[Oracle Cloud Console](https://www.oracle.com/cloud/free/)** > **Instances** > **Create Instance**.
2. Select **Canonical Ubuntu 24.04**.
3. Shape: Select **Ampere A1 Flex** (e.g., 2 OCPUs, 12GB RAM for future proofing).
4. Download the **SSH Private Key** before launching it. Wait until the instance shows **Running** and grab the **Public IP Address**.

### 2. Configure Firewalls
You must open HTTP (80) and HTTPS (443):
1. **Oracle Cloud Console:** Go to your networking subnet -> **Security Lists** -> Add Ingress Rules for `0.0.0.0/0` on TCP ports `80` and `443`.
2. **Server OS Firewall:** SSH into your instance and run:
   ```bash
   sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 80 -j ACCEPT
   sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 443 -j ACCEPT
   sudo netfilter-persistent save
   ```

### 3. Setup the Environment
SSH into the server and install Go and Caddy (Reverse Proxy):
```bash
# Install Go (arm64)
wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Install Caddy (for Auto SSL/HTTPS)
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install caddy
```

### 4. Build the Backend
```bash
git clone https://github.com/akshbswas98/kshanik_search.git
cd kshanik_search/backend

# Create your production environment file
cat <<EOF > .env
PORT=8080
PROVIDER_TIMEOUT_MS=5000
SEARCH_TIMEOUT_MS=10000
EOF

# Build the native binary
go build -o kshanik-search ./cmd/server/
```

### 5. Create Systemd Service
Keep the app running permanently across reboots:
```bash
sudo nano /etc/systemd/system/kshanik-search.service
```
Paste the following:
```ini
[Unit]
Description=Kshanik Go Search Backend
After=network.target

[Service]
User=ubuntu
Group=ubuntu
WorkingDirectory=/home/ubuntu/kshanik_search/backend
ExecStart=/home/ubuntu/kshanik_search/backend/kshanik-search
Restart=always
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```
Start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable --now kshanik-search
```

### 6. Setup Caddy for HTTPS
Link your domain (or a free one like DuckDNS) to the server's Public IP. Then edit Caddy config to point to Go:
```bash
sudo nano /etc/caddy/Caddyfile
```
```caddyfile
api.yourdomain.com {
    reverse_proxy localhost:8080
}
```
Reload config:
```bash
sudo systemctl reload caddy
```

---

## Part 3: Frontend Deployment (Netlify/Vercel)

Now that the backend is securely serving traffic at `https://api.yourdomain.com/search`, push your code changes to GitHub.

1. Connect your repository to **Netlify** or **Vercel**.
2. Deployment Settings:
   - Build Command: `npm run build`
   - Publish Directory: `dist`
3. Environment Variables:
   Add the following secret variable in your platform dashboard before deploying:
   - **Key:** `VITE_BACKEND_URL`
   - **Value:** `https://api.yourdomain.com` *(Do NOT include a trailing slash)*
4. Click **Deploy**.

**Congratulations!** Your meta-search engine is now fully deployed.
