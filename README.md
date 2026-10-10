# SEGO Server

**SEGO Server** is an open-source, self-hosted private cloud that turns an existing Windows or Linux computer into a secure personal or small-team file server.

It is designed for local ownership, low-resource machines, streaming file I/O, server-side authorization, safe sharing, auditability, emergency containment, simple installation, boot-time availability, and recoverable updates.

> **Project stage:** secure MVP under active development. Features marked as roadmap are not presented as production-complete.

## Install — one command

You do **not** need to clone the repository or download a ZIP first.

### Windows 10 / 11

Open PowerShell and paste:

```powershell
irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
```

### Ubuntu / Debian

Open Terminal as your normal user and paste:

```bash
curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
```

Do not prefix the Linux command with `sudo`; the installer asks only when elevated access is required.

The guided installer downloads SEGO from GitHub, installs supported prerequisites, chooses a resource profile, builds the dashboard/backend, creates automatic startup, creates the first admin account, enables automatic GitHub Release update checks, starts the server, and shows the dashboard address.

## Beginner commands

After installation, normal use is deliberately small:

```text
sego open          Open the dashboard
sego status        Check the server
sego start         Start SEGO
sego stop          Stop SEGO
sego restart       Restart SEGO
sego help          Show commands
```

No config path is required for these commands.

## Starts automatically after reboot

SEGO is installed as a boot-start background service/task:

- Linux: `sego.service` through systemd
- Windows: `SEGO Server` startup task running as SYSTEM in the current MVP

After the machine restarts, users do not need to open a terminal to start the server.

## Optional remote access with Cloudflare Tunnel

Run once:

```text
sego tunnel setup
```

The command installs/configures `cloudflared` as an OS service. The tunnel token is entered through a hidden prompt instead of being typed into command history. Configure the Cloudflare published application to point to:

```text
http://localhost:8088
```

After setup, the tunnel service reconnects automatically when the computer boots.

Other commands:

```text
sego tunnel status
sego tunnel remove
```

Never publish SQLite, SSH, database, or unrelated administrative ports through the SEGO hostname.

## Automatic release updates

SEGO checks the repository's **GitHub Releases** channel automatically after installation.

```text
sego version
sego update check
sego update auto
```

Automatic checks run at boot/daily depending on OS. The updater requires a matching platform package and `SHA256SUMS` asset, verifies the package checksum before installation, stops the server only after download/verification, keeps a local backup, restarts SEGO, performs a health check, and attempts rollback if the updated server fails health validation.

Source-branch changes are not silently treated as releases. If no GitHub Release has been published, the updater makes no installation change.

## Current capabilities

### Storage

- native filesystem storage roots
- per-user internal namespaces
- traversal and symlink-chain protection
- streaming upload/download
- folder creation, rename, listing and trash
- quota reservation
- expiring/revocable file share links

### Safe file previews

- image preview: JPEG, PNG, WebP, GIF, BMP, ICO
- PDF preview
- text/Markdown/JSON/CSV/log/source preview
- audio player: MP3, WAV, OGG, M4A, AAC, FLAC where browser codecs support them
- video player: MP4, WebM, OGV, MOV, M4V where browser codecs support them
- HTTP Range serving for seeking in media
- short-lived scoped preview URLs instead of login bearer tokens in media URLs
- active HTML/SVG is not admitted as inline preview content
- graceful preview-error state instead of a broken-image/player UI

### Emergency SOS

Admins have an **Emergency SOS** action for suspected compromise.

Activation immediately:

- enables global lockdown;
- revokes every active session, including the activating admin;
- blocks new remote logins and authenticated API use;
- revokes/expires all existing share links;
- invalidates temporary preview tokens;
- records the incident in the audit log.

Recovery is intentionally local-only:

```text
sego emergency unlock
```

The local command verifies an administrator password. Old share links are **not** restored after recovery.

### Audit logs

Admins can view a structured audit list and open each entry for details including:

- action name
- date/time
- username/user ID/role
- client IP
- user agent
- target/resource
- result/outcome
- structured event details

Secrets, plaintext passwords, full share tokens, and full preview tokens are not intentionally written to the audit log.

## Low-resource mode

The installer chooses `low-resource` automatically on smaller-memory systems. It can be changed later:

```text
sego mode show
sego mode low-resource
sego mode standard
```

SEGO does not claim a certified 4 GB performance target until real-hardware benchmark results are published.

## Architecture

- backend/API/CLI: Go
- metadata/auth/audit: SQLite
- file contents: native filesystem
- frontend: React + TypeScript + Vite
- Linux startup: systemd
- Windows startup: Task Scheduler in the current MVP
- optional remote access: Cloudflare Tunnel service

## Documentation

Start with [`docs/GETTING-STARTED.md`](docs/GETTING-STARTED.md).

- [`INSTALL.md`](INSTALL.md)
- [`SECURITY.md`](SECURITY.md)
- [`ARCHITECTURE.md`](ARCHITECTURE.md)
- [`CONTRIBUTING.md`](CONTRIBUTING.md)
- [`TEST-REPORT.md`](TEST-REPORT.md)
- [`docs/PREVIEWS.md`](docs/PREVIEWS.md)
- [`docs/EMERGENCY-SOS.md`](docs/EMERGENCY-SOS.md)
- [`docs/AUDIT-LOGS.md`](docs/AUDIT-LOGS.md)
- [`docs/TUNNEL.md`](docs/TUNNEL.md)
- [`docs/UPDATES.md`](docs/UPDATES.md)
- [`docs/API.md`](docs/API.md)
- [`docs/CLI.md`](docs/CLI.md)
- [`docs/THREAT-MODEL.md`](docs/THREAT-MODEL.md)

## Not yet production-complete

- MFA/TOTP
- complete user/group/folder ACL inheritance
- resumable multipart uploads
- isolated Office-document conversion preview
- full backup/restore scheduler and encrypted backup orchestration
- native signed Windows Service installer
- cryptographically signed release manifest beyond GitHub transport + SHA-256 asset validation
- complete clean-machine Windows/Linux and 4 GB hardware certification

See [`docs/ROADMAP.md`](docs/ROADMAP.md).

## Verification policy

SEGO does not claim that a test, platform target, performance target, security control, release, or rollback passed unless the relevant check was actually executed and documented.

## License

MIT License. See [`LICENSE`](LICENSE).
