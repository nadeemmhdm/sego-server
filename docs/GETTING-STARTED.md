# Getting Started

This page is intentionally written for first-time users.

## 1. Install

### Windows

Open PowerShell and paste:

```powershell
irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
```

### Ubuntu / Debian

Open Terminal and paste:

```bash
curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
```

The installer downloads SEGO from GitHub. You do not need to clone the repository or manually download a ZIP.

## 2. Create the administrator

The installation wizard asks for:

- an admin username;
- a password of at least 12 characters;
- password confirmation.

No default password is created.

## 3. Open SEGO

Use:

```text
sego open
```

Or open this address manually:

```text
http://127.0.0.1:8088
```

## 4. Everyday commands

```text
sego open       Open dashboard
sego status     Check server
sego start      Start server
sego stop       Stop server
sego restart    Restart server
sego help       Help
```

## 5. Upload and preview files

Sign in, open **Files**, then choose **Upload**. Click a supported file to preview it. Supported browser codecs can play audio and video directly with seeking.

## 6. Security emergency

Admins can open **Security** and use **Emergency SOS** if an active compromise is suspected. It immediately blocks access and revokes sessions/links.

To recover on the server computer:

```text
sego emergency unlock
```

## 7. If something does not work

Run:

```text
sego status
sego doctor
```

Then see [`TROUBLESHOOTING.md`](TROUBLESHOOTING.md).
