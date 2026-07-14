# Final Deployment Guide: Kshanik Search Engine

Deploy the React frontend to Netlify and the Go backend to an Oracle Cloud
Always Free Ampere A1 VM.

## Why Oracle Cloud

Oracle documents an Always Free Ampere A1 allowance of 3,000 OCPU-hours and
18,000 GB-hours per month, equivalent to 4 OCPUs and 24 GB RAM. The project only
needs one small `1` OCPU / `6 GB` VM initially. You retain SSH access and enough
capacity to add services later.

Fly.io is easier to start with, but it is not the free-tier choice: its current
documentation says there is no free account or free tier for new users.

OCI can temporarily run out of Always Free A1 capacity in your home region and
Oracle documents reclamation of inactive Always Free compute instances. This is
a strong hobby-project deployment with control, not a production SLA.

## Backend

The backend directory contains a Docker Compose deployment:

- `backend`: the Go API, private inside the Docker network
- `caddy`: public reverse proxy on ports `80` and `443` with automatic HTTPS
- `deploy.sh`: Docker installation, systemd registration, and deployment

Follow [backend/DEPLOYMENT.md](backend/DEPLOYMENT.md) for the OCI VM, DuckDNS,
firewall, and deployment commands.

## Frontend

The frontend already reads:

```javascript
const API_BASE = import.meta.env.VITE_BACKEND_URL || '';
```

In Netlify, set:

```text
VITE_BACKEND_URL=https://YOUR-DUCKDNS-SUBDOMAIN.duckdns.org
```

Use:

| Setting | Value |
| --- | --- |
| Build command | `npm run build` |
| Publish directory | `dist` |

Trigger a Netlify deployment after changing the environment variable.

## Verification

```bash
curl https://YOUR-DUCKDNS-SUBDOMAIN.duckdns.org/health
curl "https://YOUR-DUCKDNS-SUBDOMAIN.duckdns.org/search?q=test&limit=2"
```

Open the Netlify site and confirm in browser DevTools that requests go to the
DuckDNS HTTPS hostname without CORS errors.
