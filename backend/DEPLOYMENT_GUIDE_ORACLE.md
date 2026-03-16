# Deploying Kshanik Search Backend on Oracle Cloud (Always Free ARM Instance)

This guide walks you through deploying your highly-concurrent Go meta-search backend onto an Oracle Cloud Infrastructure (OCI) ARM instance. 

Oracle provides a very generous "**Always Free**" tier: Up to 4 ARM Ampere A1 Compute instances with 24GB of RAM. Since Go cross-compiles perfectly to `linux/arm64`, this is the ideal setup for your future Web Crawler and Indexer running natively under Systemd.

---

## Phase 1: Create Your OCI Free Instance

1.  **Sign Up:** Create a Free Tier account at [Oracle Cloud](https://www.oracle.com/cloud/free/).
2.  **Launch a VM:** In the OCI Console, go to **Compute > Instances** and click **Create Instance**.
3.  **Configure the Instance:**
    *   **Image:** Choose **Canonical Ubuntu 24.04** (or 22.04).
    *   **Shape:** Click "Change Shape". Select **Virtual Machine** -> **Ampere** -> **VM.Standard.A1.Flex**. 
    *   *Tip:* You can assign up to 4 OCPUs and 24GB RAM for free. Allocate what you feel comfortable with (even 1 OCPU/6GB RAM is overkill for this backend).
4.  **Networking:** Ensure the instance gets a Public IP address (this is checked by default).
5.  **SSH Keys:** Under "Add SSH keys", download the **Private Key** (and the Public Key) they generate for you. *Do not lose this private key.*
6.  **Launch:** Click "Create". Wait about 1-2 minutes for the instance status to change to **RUNNING**, and note the **Public IP Address**.

---

## Phase 2: Open Firewall Ports (Oracle Cloud + Ubuntu)

Oracle has strict firewalls at both the cloud-network level and the OS level. We need to open HTTP (80) and HTTPS (443) to expose the API to the Internet.

### 2A: Oracle Cloud VCN Dashboard
1. On the Instance Details page, click the associated **Subnet** link.
2. Click on the **Security List** for that subnet.
3. Add two **Ingress Rules**:
   - **CIDR:** `0.0.0.0/0` | **Protocol:** `TCP` | **Destination Port:** `80` (Caddy/HTTP)
   - **CIDR:** `0.0.0.0/0` | **Protocol:** `TCP` | **Destination Port:** `443` (Caddy/HTTPS)

### 2B: Ubuntu `iptables` (SSH into Server)
Open a terminal locally and SSH into your new server using the private key you downloaded:

```bash
# Set strict permissions on the keyfile
chmod 400 your_private_key.key

# Connect to the server
ssh -i your_private_key.key ubuntu@YOUR_PUBLIC_IP
```

Now configure the OS firewall on the server to allow traffic:

```bash
# Allow web traffic
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 80 -j ACCEPT
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 443 -j ACCEPT
sudo netfilter-persistent save
```

---

## Phase 3: Setup the Server Environment

You are now logged into the Ubuntu ARM server.

### 1. Install Go
Download and install the latest `arm64` version of Go.

```bash
cd ~
wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Verify installation
go version
# It should print: go version go1.22.0 linux/arm64
```

### 2. Install Reverse Proxy (Caddy)
We will use **Caddy** instead of Nginx. Caddy automatically provisions and renews Let's Encrypt SSL certificates (which is required by React/Netlify — you cannot make an insecure `http://` request from a secure `https://` Netlify site).

```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install caddy
```

---

## Phase 4: Build and Deploy Your Application

**Locally (on your machine):** Make sure your `go-backend-refactor` branch is pushed to GitHub.

**On the Oracle Server:**
```bash
# Clone your repository
git clone https://github.com/akshbswas98/kshanik_search.git
cd kshanik_search/backend

# Create your production .env file
cat <<EOF > .env
PORT=8080
PROVIDER_TIMEOUT_MS=5000
SEARCH_TIMEOUT_MS=10000
GITHUB_TOKEN=your_token_if_you_have_one
EOF

# Build the executable natively for ARM
go build -o kshanik-search ./cmd/server/
```

---

## Phase 5: Run as a Persistent Systemd Service

We don't want the server to stop when we close the SSH window, nor do we want to manually restart it if the instance reboots. We create a `systemd` service for that.

```bash
sudo nano /etc/systemd/system/kshanik-search.service
```

Paste the following configuration:

```ini
[Unit]
Description=Kshanik Go Search Backend
After=network.target

[Service]
# User and Group
User=ubuntu
Group=ubuntu

# Set the working directory (where the .env file is)
WorkingDirectory=/home/ubuntu/kshanik_search/backend

# Provide the full path to the compiled binary
ExecStart=/home/ubuntu/kshanik_search/backend/kshanik-search

# Restart policies
Restart=always
RestartSec=5

# Limit open files (useful for future crawler)
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

Enable and start your backend:

```bash
sudo systemctl daemon-reload
sudo systemctl enable kshanik-search
sudo systemctl start kshanik-search

# Check the logs to ensure it's running cleanly on port 8080!
sudo journalctl -u kshanik-search -f
```

---

## Phase 6: Configure HTTPS via Caddy

Netlify requires the backend to serve traffic over `https://`. You need a domain name (or a free subdomain like DuckDNS).

Point your domain's DNS `A` record (e.g., `api.kshaniksearch.com`) to your **Oracle Server's Public IP Address**.

Next, open the Caddy configuration file:

```bash
sudo nano /etc/caddy/Caddyfile
```

Replace everything in it with:

```caddyfile
api.yourdomain.com {
    reverse_proxy localhost:8080
}
```

Reload Caddy to apply the changes and automatically request an SSL certificate:

```bash
sudo systemctl reload caddy
```

Your API is now live, secure, and permanently running at `https://api.yourdomain.com/search?q=...`!

---

## Phase 7: Update React (Netlify) & Deploy

Now that your Go backend is deployed and serving results securely, update the frontend configuration to point to it.

1. **In your local code (`src/components/Results.jsx`),** make sure it hits the full domain (or use Vite environment variables):
   ```javascript
   const API_BASE = import.meta.env.VITE_BACKEND_URL || ''; 
   const response = await fetch(`${API_BASE}/api/search?${params.toString()}`);
   // Important: Keep the /search path the same as what the Go Handler expects.
   ```
2. **In Netlify Settings:** Go to *Project Settings > Environment Variables > Environment variables* and add a new secret:
   - **Key:** `VITE_BACKEND_URL`
   - **Value:** `https://api.yourdomain.com` (no trailing slash)
   *(Note: The `vite.config.js` proxy array is completely ignored by Netlify in production; Netlify builds static files that hit the URL you define).*
3. Add, commit, and push your changes to the `main` branch. Netlify will auto-deploy the site, and it will query your insanely fast Oracle Cloud Go instance!
