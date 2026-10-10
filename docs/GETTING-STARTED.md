# Getting Started

This guide is for first-time users. You only need one install command and a few short SEGO commands.

## Step 1 — Install

### Windows

Open PowerShell and paste:

```powershell
irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
```

Approve the Administrator prompt.

### Ubuntu / Debian

Open Terminal and paste:

```bash
curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
```

Do not manually download a ZIP or clone the repository. The installer retrieves SEGO directly from GitHub.

## Step 2 — Create the administrator

When asked, enter an admin username and a password of at least 12 characters. SEGO does not create a default password.

## Step 3 — Open SEGO

After installation simply run:

```text
sego
```

This starts SEGO if necessary and opens the dashboard. The local address is `http://127.0.0.1:8088`.

## Commands worth remembering

```text
sego             Open dashboard
sego status      Check server
sego update      Update SEGO
sego remote      Set up remote access
sego unlock      Recover after Emergency SOS
sego help        Show all commands
```

Service controls are also available: `sego start`, `sego stop`, and `sego restart`.

## Files and previews

Sign in, open **Files**, and choose **Upload**. Click a supported file to preview it. Supported audio/video files use the built-in browser player and allow seeking when the browser codec supports the file.

## Emergency SOS

If a security incident is suspected, an admin can open **Security** and activate **Emergency SOS**. It blocks access and revokes active sessions and links.

Recovery must be performed locally on the server computer:

```text
sego unlock
```

## Troubleshooting

Run:

```text
sego status
sego doctor
```

Then see [`TROUBLESHOOTING.md`](TROUBLESHOOTING.md).
