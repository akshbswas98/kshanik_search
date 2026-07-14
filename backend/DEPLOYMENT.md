# Deploy the Kshanik Backend on Oracle Cloud Always Free

This guide deploys the Kshanik Search backend to the internet. Follow the steps
in order. Do not skip a step.

The result will look like this:

```text
Browser
  |
  | HTTPS request
  v
https://YOUR-NAME.duckdns.org
  |
  v
Oracle Cloud VM
  |
  +-- Caddy container: receives HTTPS requests
  |
  +-- Go backend container: runs the search API privately
```

## What you are creating

| Item | Purpose | Cost |
| --- | --- | --- |
| Oracle Cloud Ampere A1 VM | A small computer running on the internet | Always Free allowance |
| Docker Compose | Starts and restarts the backend containers | Free |
| Caddy | Adds HTTPS automatically | Free |
| DuckDNS subdomain | Gives the backend an easy-to-use internet address | Free |

Oracle documents an Always Free Ampere A1 allowance of 3,000 OCPU-hours and
18,000 GB-hours per month. This equals 4 OCPUs and 24 GB RAM. This guide uses
only 1 OCPU and 6 GB RAM.

Important limitations:

- Create Always Free resources only in your Oracle tenancy's home region.
- Oracle may temporarily have no free Ampere A1 capacity. If that happens, try
  another availability domain or try again later.
- Oracle documents reclamation of inactive Always Free compute instances.
- This is a strong hobby-project deployment, not hosting with a production SLA.

## Before you start

You need:

1. An [Oracle Cloud account](https://www.oracle.com/cloud/free/).
2. A [DuckDNS account](https://www.duckdns.org/). You can sign in with an
   existing Google, GitHub, Reddit, or other supported account.
3. Your frontend Netlify address, such as:

   ```text
   https://kshaniksearch.netlify.app
   ```

4. A computer with a terminal:

   - On Windows, open **PowerShell**.
   - On macOS or Linux, open **Terminal**.

Whenever this guide shows `YOUR-NAME`, replace it with your own value. Do not
type the words `YOUR-NAME`.

## Step 1: Create an SSH key

An SSH key is a secure digital key. You use it to enter your Oracle VM.

### Windows PowerShell

Run:

```powershell
ssh-keygen -t ed25519 -C "kshanik-oracle"
```

When PowerShell asks where to save the key, press **Enter**. When it asks for a
passphrase, either enter a passphrase you will remember or press **Enter**
twice to continue without one.

Show your public key:

```powershell
Get-Content $HOME\.ssh\id_ed25519.pub
```

### macOS or Linux

Run:

```bash
ssh-keygen -t ed25519 -C "kshanik-oracle"
cat ~/.ssh/id_ed25519.pub
```

Copy the full line that appears. It starts with `ssh-ed25519`. This is the
public key. It is safe to paste into Oracle Cloud.

Keep the private key secret. The private key is the file without `.pub` at the
end. Never paste it into a website, chat, GitHub repository, or `.env` file.

## Step 2: Create the Oracle VM

1. Open the [Oracle Cloud Console](https://cloud.oracle.com/).
2. Confirm the region shown near the top-right corner is your home region.
3. Open the top-left menu.
4. Click **Compute**.
5. Click **Instances**.
6. Click **Create instance**.
7. Enter this name:

   ```text
   kshanik-search
   ```

8. Under **Image and shape**, click **Edit**.
9. For the image, click **Change image**.
10. Select **Canonical Ubuntu** and choose **Ubuntu 24.04**.
11. Click **Select image**.
12. For the shape, click **Change shape**.
13. Choose **Ampere**.
14. Select:

   ```text
   VM.Standard.A1.Flex
   ```

15. Set:

   | Setting | Value |
   | --- | --- |
   | OCPUs | `1` |
   | Memory | `6 GB` |

16. Click **Select shape**.
17. Fill in the **Networking** section using the detailed instructions below.
18. Fill in the **Add SSH keys** section using the detailed instructions below.
19. Click **Create**.

### Networking choices

The VM needs one network connection called the **Primary VNIC**. Your Oracle
account already has a VCN named `vcn-20260426-1727` and a regional subnet named
`subnet-20260426-1727`. Reuse them to keep the setup simple.

Choose these values:

| Screen option | What to choose |
| --- | --- |
| VNIC name | Leave empty |
| Primary network | **Select existing virtual cloud network** |
| Virtual cloud network compartment | `biswaakashab (root)` |
| Virtual cloud network | `vcn-20260426-1727` |
| Subnet | **Select existing subnet** |
| Subnet compartment | `biswaakashab (root)` |
| Subnet | `subnet-20260426-1727 (regional)` |
| Private IPv4 address | **Automatically assign private IPv4 address** |
| Public IPv4 address | Enable **Automatically assign public IPv4 address** |
| IPv6 address | Leave disabled |

Why there are two IP addresses:

- The automatically assigned **private IPv4 address** is used inside Oracle's
  network. It will look like `10.0.0.x`.
- The automatically assigned **public IPv4 address** lets your computer,
  DuckDNS, and website reach the VM through the internet.

Do not choose:

- **Create new virtual cloud network**
- **Specify OCID**
- **Create new public subnet**
- **Manually assign private IPv4 address**
- **Provide existing private IPv4 OCID**

Those options are useful for advanced setups but add unnecessary work here.

### Advanced networking choices

Expand **Advanced options** only if Oracle shows them. Choose:

| Screen option | What to choose |
| --- | --- |
| Use network security groups to control traffic | Leave disabled |
| DNS record | **Assign a private DNS record** |
| Hostname | Leave the automatically generated value |
| Launch options | **Let Oracle Cloud Infrastructure choose the best networking type** |

You do not need a network security group for this first deployment. The
security-list rules added in Step 3 control access.

The private DNS record is only for Oracle's internal network. Your public
internet hostname is created later with DuckDNS.

### Confirm that the existing subnet is public

Before creating the VM, confirm that the selected existing subnet can reach the
internet:

1. Open the top-left Oracle Cloud menu in a second browser tab.
2. Click **Networking**.
3. Click **Virtual cloud networks**.
4. Click `vcn-20260426-1727`.
5. Click **Subnets**.
6. Click `subnet-20260426-1727`.
7. Confirm that **Prohibit public IP on VNIC** is `No`.
8. Click the route table linked from the subnet page.
9. Confirm that it has a route rule with:

   | Field | Value |
   | --- | --- |
   | Destination | `0.0.0.0/0` |
   | Target type | `Internet Gateway` |

If both checks pass, return to the Create Instance browser tab and continue.

If **Prohibit public IP on VNIC** is `Yes`, do not use that subnet. Return to
the Create Instance tab, select **Create new public subnet**, give it a name
such as `kshanik-public-subnet`, and keep Oracle's default values.

### SSH key choices

If you completed Step 1 and copied your public key:

1. Select **Paste public key**.
2. Paste the line beginning with `ssh-ed25519`.

Do not paste the private key.

If you skipped Step 1, choose **Generate a key pair for me** instead. Then:

1. Click **Download private key**.
2. Store the downloaded file somewhere safe.
3. Click **Download public key** if you want a backup.

Oracle shows the generated private key only once. If you lose it, you cannot
use it to connect later.

Wait until the instance status changes from `Provisioning` to `Running`.

Find and copy the **Public IPv4 address** from the instance page. It looks like:

```text
123.45.67.89
```

Write it down. This guide calls it `YOUR-ORACLE-IP`.

## Step 3: Open the Oracle Cloud firewall

Your backend needs:

| Port | Used for |
| --- | --- |
| `22` | SSH terminal access |
| `80` | Initial HTTP connection and HTTPS certificate setup |
| `443` | Secure HTTPS traffic |

From the VM instance page:

1. Scroll to **Primary VNIC** and click the subnet name.
2. Click the security list attached to the subnet. It is often named
   **Default Security List for ...**.
3. Click **Add ingress rules**.
4. Add an SSH rule:

   | Field | Value |
   | --- | --- |
   | Source CIDR | `0.0.0.0/0` |
   | IP Protocol | `TCP` |
   | Destination Port Range | `22` |
   | Description | `SSH` |

5. Add an HTTP rule:

   | Field | Value |
   | --- | --- |
   | Source CIDR | `0.0.0.0/0` |
   | IP Protocol | `TCP` |
   | Destination Port Range | `80` |
   | Description | `HTTP for Caddy` |

6. Add an HTTPS rule:

   | Field | Value |
   | --- | --- |
   | Source CIDR | `0.0.0.0/0` |
   | IP Protocol | `TCP` |
   | Destination Port Range | `443` |
   | Description | `HTTPS for Caddy` |

7. Click **Add ingress rules**.

For a first deployment, allowing SSH from `0.0.0.0/0` is simpler. After the
deployment works, replace the SSH rule source with your own public IP followed
by `/32`, such as `123.45.67.89/32`.

Do not add a public rule for port `8080`. The Go API stays private behind
Caddy.

## Step 4: Create the free DuckDNS address

An IP address is difficult to remember. DuckDNS gives it a name and lets Caddy
create an HTTPS certificate.

1. Open [DuckDNS](https://www.duckdns.org/).
2. Sign in.
3. Find **domains**.
4. Enter a unique subdomain. For example:

   ```text
   api-kshanik-yourname
   ```

5. Click **add domain**.
6. Find your new subdomain in the list.
7. In the **current ip** field, enter your Oracle VM public IPv4 address from
   Step 2.
8. Click **update ip**.

Your complete backend hostname is now:

```text
api-kshanik-yourname.duckdns.org
```

This guide calls it `YOUR-NAME.duckdns.org`.

Wait one or two minutes. Then check that the name points to your VM.

### Windows PowerShell

```powershell
nslookup YOUR-NAME.duckdns.org
```

### macOS or Linux

```bash
nslookup YOUR-NAME.duckdns.org
```

Look for your Oracle VM public IPv4 address in the output. Continue only after
the correct IP appears.

## Step 5: Connect to the Oracle VM

Open PowerShell or Terminal on your computer.

Run this command after replacing `YOUR-ORACLE-IP`:

```bash
ssh ubuntu@YOUR-ORACLE-IP
```

The first time you connect, your terminal may ask:

```text
Are you sure you want to continue connecting (yes/no/[fingerprint])?
```

Type:

```text
yes
```

Press **Enter**. You are now inside the Oracle VM. Commands from the next steps
must run inside this SSH terminal.

## Step 6: Allow web traffic inside the Ubuntu VM

Oracle Cloud has an outer firewall configured in Step 3. The Ubuntu VM image
can also have an inner firewall. Run:

```bash
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 80 -j ACCEPT
sudo iptables -I INPUT 6 -m state --state NEW -p tcp --dport 443 -j ACCEPT
```

Install the tool that saves these rules after a reboot:

```bash
sudo apt-get update
sudo apt-get install -y iptables-persistent
```

During installation, choose **Yes** when asked to save the current IPv4 rules.
Choose **Yes** for IPv6 rules too.

Save once more:

```bash
sudo netfilter-persistent save
```

## Step 7: Download the project

Still inside the SSH terminal, install Git:

```bash
sudo apt-get install -y git
```

Download the repository:

```bash
git clone https://github.com/akshbswas98/kshanik_search.git
```

Move into the backend folder:

```bash
cd kshanik_search/backend
```

Make the deployment script executable:

```bash
chmod +x deploy.sh
```

## Step 8: Run the setup script once

Run:

```bash
./deploy.sh
```

The script installs Docker if it is missing. The first run also creates a file
named `.env` and stops so that you can fill in your settings.

The final lines should say:

```text
Creating .env from .env.example...
Edit .env with production values, then run ./deploy.sh again.
```

## Step 9: Fill in the production settings

Open the `.env` file:

```bash
nano .env
```

Use the arrow keys to move around. Change the file so these lines contain your
real values:

```dotenv
ENV=production
API_DOMAIN=YOUR-NAME.duckdns.org
ACME_EMAIL=YOUR-EMAIL@example.com
ALLOWED_ORIGINS=https://YOUR-NETLIFY-SITE.netlify.app
```

Example:

```dotenv
ENV=production
API_DOMAIN=api-kshanik-akash.duckdns.org
ACME_EMAIL=akash@example.com
ALLOWED_ORIGINS=https://kshaniksearch.netlify.app
```

Rules:

- Do not include `https://` in `API_DOMAIN`.
- Do include `https://` in `ALLOWED_ORIGINS`.
- Do not add a `/` at the end of either value.
- `ACME_EMAIL` should be an email address you can access.
- Keep `PORT=8080`.

The other default settings can remain unchanged.

Optional: add a GitHub token to improve GitHub search rate limits:

```dotenv
GITHUB_TOKEN=
```

Leave it empty for the first deployment if you do not already have a token.

Save the file:

1. Press **Ctrl+O**.
2. Press **Enter**.
3. Press **Ctrl+X**.

## Step 10: Start the backend

Run:

```bash
./deploy.sh
```

This may take a few minutes. Docker downloads the required images and builds
the Go backend.

When it finishes, check the containers:

```bash
sudo docker compose ps
```

Look for two containers:

| Container | Expected state |
| --- | --- |
| `kshanik-backend` | `Up` or `healthy` |
| `kshanik-proxy` | `Up` |

## Step 11: Test the backend

Test the health endpoint:

```bash
curl https://YOUR-NAME.duckdns.org/health
```

Then test a search:

```bash
curl "https://YOUR-NAME.duckdns.org/search?q=test&limit=2"
```

If both commands return a response without an error, your backend is live on
the internet.

## Step 12: Connect the Netlify frontend

The backend is deployed. Now tell the frontend where to find it.

1. Open [Netlify](https://app.netlify.com/).
2. Open your Kshanik frontend site.
3. Click **Site configuration**.
4. Click **Environment variables**.
5. Click **Add a variable**.
6. Enter:

   | Field | Value |
   | --- | --- |
   | Key | `VITE_BACKEND_URL` |
   | Value | `https://YOUR-NAME.duckdns.org` |

7. Save the variable.
8. Open **Deploys**.
9. Click **Trigger deploy**.
10. Click **Deploy site**.

Wait for Netlify to finish. Open the frontend website and perform a search.

## Step 13: Confirm the complete deployment

Open the Netlify frontend in your browser.

1. Search for something simple, such as `golang`.
2. Confirm that results appear.
3. If results do not appear, press **F12** to open browser developer tools.
4. Open the **Network** tab.
5. Search again.
6. Confirm that the request address starts with:

   ```text
   https://YOUR-NAME.duckdns.org/search
   ```

## Updating the backend later

SSH into the VM:

```bash
ssh ubuntu@YOUR-ORACLE-IP
```

Move into the backend folder:

```bash
cd kshanik_search/backend
```

Download repository changes and redeploy:

```bash
git pull
./deploy.sh
```

## Useful commands

Run these commands inside `~/kshanik_search/backend` on the Oracle VM.

| Goal | Command |
| --- | --- |
| Show container status | `sudo docker compose ps` |
| Watch all container logs | `sudo docker compose logs -f` |
| Watch backend logs only | `sudo docker compose logs -f backend` |
| Watch Caddy HTTPS logs only | `sudo docker compose logs -f caddy` |
| Restart the deployment | `./deploy.sh` |
| Watch boot service logs | `sudo journalctl -u kshanik-stack -f` |

Press **Ctrl+C** to stop watching logs. This does not stop the backend.

## Troubleshooting

### `Out of host capacity` appears while creating the VM

Oracle does not currently have a free Ampere A1 machine available in that
availability domain. Choose a different availability domain if your region has
one, or retry later.

### `ssh: connect to host ... port 22: Connection timed out`

Check:

1. The VM status is `Running`.
2. You used the correct public IPv4 address.
3. The Oracle security list has the TCP port `22` ingress rule from Step 3.

### DuckDNS does not show the Oracle IP with `nslookup`

Open DuckDNS again. Confirm the **current ip** field contains the Oracle public
IPv4 address. Click **update ip**, wait one or two minutes, and retry.

### `curl` returns a connection timeout

Check that the Oracle security list has ingress rules for TCP ports `80` and
`443`. Then rerun the Ubuntu firewall commands from Step 6.

### `curl` returns an HTTPS or certificate error

Check:

1. `API_DOMAIN` contains only the hostname, without `https://`.
2. DuckDNS points to the Oracle public IPv4 address.
3. TCP ports `80` and `443` are open.

Then inspect Caddy logs:

```bash
sudo docker compose logs -f caddy
```

### The Netlify site shows a CORS error

Open `.env`:

```bash
nano .env
```

Confirm `ALLOWED_ORIGINS` exactly matches the frontend address:

```dotenv
ALLOWED_ORIGINS=https://YOUR-NETLIFY-SITE.netlify.app
```

Do not add a trailing `/`. Save the file and run:

```bash
./deploy.sh
```

### A container is not running

Show the logs:

```bash
sudo docker compose logs
```

The last lines usually explain the problem. After fixing the setting, run:

```bash
./deploy.sh
```
