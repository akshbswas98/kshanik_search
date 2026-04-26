# 📝 Kshanik Search - Deployment Session Log
**Date:** Sunday, 26 April 2026
**Status:** In-Progress (Server Setup Phase)

---

## 🏗️ Local Development Status
- **Backend:** Go server (v1.26.0) running on `localhost:8080`.
- **Frontend:** Vite React server running on `localhost:3000`.
- **Git Branch:** `stage-kshanik-search` (Pushed to GitHub with latest changes).
- **Builds:** 
  - Frontend `dist/` is ready.
  - Backend Windows `.exe` and Linux `arm64` binaries are built in `backend/bin/`.

---

## ☁️ Oracle Cloud Instance Details
- **Instance Name:** `kshanik-search-prod-01`
- **Public IP:** `80.225.228.75`
- **Shape:** `VM.Standard.E2.1.Micro` (Always Free)
- **OS:** Oracle Linux 9
- **Username:** `opc`
- **SSH Key Path:** `C:\Users\Akash\Downloads\ksk-search-keys-oracle\final-2\ssh-key-2026-04-26.key`

### 🔑 SSH Login Command
```powershell
ssh -i "C:\Users\Akash\Downloads\ksk-search-keys-oracle\final-2\ssh-key-2026-04-26.key" opc@80.225.228.75
```

---

## 🛠️ Server Progress & Next Steps
We were in the middle of setting up the server. The `dnf install` command for Docker appeared to be processing slowly due to the 1GB RAM limit on the Micro instance.

### Planned "Magic Setup" Script:
Once logged back in, these are the commands to run (one by one is safer for memory):

1. **Install Docker & Git:**
   ```bash
   sudo dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
   sudo dnf install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin git
   sudo systemctl enable --now docker
   sudo usermod -aG docker opc
   ```

2. **Configure Firewall:**
   ```bash
   sudo firewall-cmd --permanent --add-service=http
   sudo firewall-cmd --permanent --add-service=https
   sudo firewall-cmd --reload
   ```

3. **Deploy Code:**
   ```bash
   git clone -b stage-kshanik-search https://github.com/akshbswas98/kshanik_search.git
   cd kshanik_search/backend
   cp .env.example .env
   sudo docker compose up -d --build
   ```

---

## 🚌 Note for Resuming
If the `dnf` command is still stuck when you return, run `sudo dnf clean all` and try installing `git` first to verify the package manager is responding. 

Safe travels on your bus! You are 90% of the way there.
