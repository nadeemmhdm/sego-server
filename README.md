# SEGO Server

**SEGO Server** is an open-source, self-hosted private cloud that turns a Windows or Linux computer into a secure personal or small-team file server.

SEGO focuses on simple installation, local ownership, low resource use, streaming file I/O, server-side authorization, safe sharing, auditability, emergency containment, automatic startup and recoverable updates.

> **Project stage:** secure MVP under active development. Roadmap features are not presented as production-complete.

## Install in one command

You do **not** need to clone GitHub, download a ZIP, install the source manually, or type build commands.

### Windows 10 / 11

Open **PowerShell** and paste one command:

```powershell
irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
```

Approve the Windows Administrator prompt when it appears. SEGO downloads itself directly from GitHub and guides you through the first admin account.

### Ubuntu / Debian

Open **Terminal** as your normal user and paste one command:

```bash
curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
```

Do not add `sudo` in front. The installer asks for sudo only when the operating system requires it.

### What the installer handles automatically

It downloads the current GitHub source snapshot, installs supported prerequisites, selects a low-resource or standard profile, resolves Go dependencies, builds the web dashboard and binaries, creates protected data/config folders, enables boot startup and update checks, creates the first Super Admin, starts SEGO and performs health checks.

## Beginner use

After installation, the shortest workflow is:

```text
sego             Open SEGO in your browser
sego status      Check whether SEGO is running
sego update      Install an available verified release update
sego remote      Set up Cloudflare Tunnel remote access
sego unlock      Recover locally after Emergency SOS
sego help        Show the full command list
```

Useful service commands remain simple:

```text
sego start
sego stop
sego restart
```

The longer administrator commands continue to work for compatibility, but beginners do not need to memorize them.

## Starts automatically

SEGO starts with the computer after installation:

- Linux: systemd `sego.service`
- Windows: `SEGO Server` startup task in the current MVP

Normally a user only needs to run `sego` to open the dashboard.

## Main capabilities

### Files and storage

- native filesystem storage roots
- per-user internal namespaces
- traversal and symlink-chain protection
- streaming upload/download
- folders, rename, listing and trash
- quota reservation
- expiring/revocable file share links

### Safe previews and media player

- images: JPEG, PNG, WebP, GIF, BMP, ICO
- PDF
- text, Markdown, JSON, CSV, logs and source files
- audio player for common browser-supported formats
- video player for MP4, WebM, OGV, MOV and M4V when the browser supports the codec
- HTTP Range requests for media seeking
- short-lived scoped preview URLs instead of exposing the login bearer token
- HTML/SVG active content is not rendered as an inline preview
- graceful preview-error UI instead of broken image/player states

### Emergency SOS

Admins can activate **Emergency SOS** during a suspected compromise. One action enables global lockdown, revokes all sessions, blocks remote login/API access, revokes and expires share links, invalidates preview tokens and writes the event to the audit log.

Recovery is intentionally local-only. On the server computer run:

```text
sego unlock
```

The admin password is verified before lockdown is removed. Revoked share links remain revoked.

### Audit logs

The admin Audit Logs view shows action, date/time, user identity and role when available, client IP, user agent, target/resource, outcome and structured details. Secrets and full tokens are not intentionally logged.

## Performance modes

SEGO selects a suitable mode during installation. Manual aliases are available:

```text
sego low        Use Low Resource Mode
sego normal     Use Standard Mode
```

SEGO does not claim certified 4 GB performance until real-hardware benchmarks are published.

## Documentation

Start here: [`docs/GETTING-STARTED.md`](docs/GETTING-STARTED.md)

- [`INSTALL.md`](INSTALL.md)
- [`SECURITY.md`](SECURITY.md)
- [`ARCHITECTURE.md`](ARCHITECTURE.md)
- [`docs/CLI.md`](docs/CLI.md)
- [`docs/PREVIEWS.md`](docs/PREVIEWS.md)
- [`docs/EMERGENCY-SOS.md`](docs/EMERGENCY-SOS.md)
- [`docs/AUDIT-LOGS.md`](docs/AUDIT-LOGS.md)
- [`docs/TUNNEL.md`](docs/TUNNEL.md)
- [`docs/UPDATES.md`](docs/UPDATES.md)
- [`docs/THREAT-MODEL.md`](docs/THREAT-MODEL.md)

## Not yet production-complete

MFA/TOTP, complete user/group/folder ACL inheritance, resumable multipart uploads, isolated Office conversion previews, full backup/restore orchestration, a native signed Windows Service installer, a cryptographically signed release manifest beyond GitHub transport plus SHA-256 asset validation, and complete clean-machine/4 GB hardware certification remain roadmap work.

See [`docs/ROADMAP.md`](docs/ROADMAP.md).

## Verification policy

SEGO does not claim a test, platform target, performance target, security control, release or rollback passed unless the relevant check was actually executed and documented.

## License

MIT License. See [`LICENSE`](LICENSE).
