# SEGO Server

**SEGO Server** is an open-source, self-hosted private cloud storage platform designed to turn an existing Windows or Linux computer into a secure personal or small-team file server.

The project prioritizes data safety, explicit authorization, streaming I/O, low resource consumption, auditable sharing, and recoverable operations over decorative features.

> **Project stage:** secure MVP foundation. Features marked as planned are not presented as production-complete.

## Why SEGO Server

SEGO is intended for users who want to keep files on hardware they control while still getting a modern browser-based file experience. The architecture uses a Go API/CLI, SQLite metadata, the native filesystem, and a React/TypeScript frontend.

### Current MVP capabilities

- Go API server and `sego` command-line interface
- SQLite-backed users, sessions, file metadata, quotas, and share records
- Argon2id password hashing
- Hashed bearer-session and share tokens
- Server-side authorization checks
- Per-user filesystem namespaces
- Storage-root confinement and traversal protection
- Streaming uploads and downloads
- Folder creation, rename, listing, and trash operations
- Conservative quota reservation for concurrent uploads
- Revocable and expiring file-share links
- Resource diagnostics endpoint
- React + TypeScript + Vite dashboard foundation
- Linux systemd installation path
- Windows source installation path
- Low-resource and standard runtime modes
- Automated Go/web CI configuration
- Security, architecture, operations, and contributor documentation

### Not yet production-complete

The following remain roadmap items and must not be treated as completed merely because interfaces or documentation exist:

- MFA/TOTP
- complete user/group/folder ACL inheritance
- resumable multipart upload protocol
- full preview pipeline and Office preview isolation
- Cloudflare Tunnel lifecycle automation
- signed update/rollback pipeline
- scheduled backup engine and restore orchestration
- native Windows Service integration/signing
- complete audit/security monitoring
- real-hardware 4 GB benchmark certification

See [`docs/ROADMAP.md`](docs/ROADMAP.md) for the staged plan.

## Supported platform baseline

| Platform | Status | Notes |
|---|---|---|
| Ubuntu Server | Initial target | systemd installer path |
| Debian | Initial target | systemd installer path |
| Windows 10 | MVP source support | native service integration pending |
| Windows 11 | MVP source support | native service integration pending |
| amd64/x86_64 | Target | primary architecture |
| arm64/aarch64 | Target | requires dependency/build validation |

Detailed support policy: [`docs/SUPPORTED-OS.md`](docs/SUPPORTED-OS.md).

## Installation

### Linux — one command from a source checkout

```bash
sudo ./scripts/install.sh --profile low-resource --yes
```

Standard mode:

```bash
sudo ./scripts/install.sh --profile standard --yes
```

### Windows — one command from an elevated PowerShell source checkout

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Profile low-resource -Yes
```

Standard mode:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Profile standard -Yes
```

The source installers require Go 1.23+ because they build the binaries locally. A future signed release installer can remove that source-build requirement once the release/signing pipeline is implemented and verified.

Full instructions: [`INSTALL.md`](INSTALL.md).

## Runtime mode

SEGO can be switched after installation without reinstalling.

```bash
sego --config /etc/sego/sego.json mode show
sego --config /etc/sego/sego.json mode set low-resource
sego --config /etc/sego/sego.json mode set standard
```

`low-resource` reduces upload concurrency and is the recommended default for older 4 GB-class systems. The project does **not** claim that the 4 GB target is validated until measured benchmark results are published.

## First administrator

Linux:

```bash
sudo /opt/sego/bin/sego --config /etc/sego/sego.json admin create
```

Windows:

```powershell
& "$env:ProgramFiles\SEGO\bin\sego.exe" --config "$env:ProgramData\SEGO\sego.json" admin create
```

Passwords are read without intentionally placing them in command-line arguments or logs.

## Development

Backend:

```bash
go run ./cmd/server --config ./sego.dev.json
```

Create a development admin:

```bash
go run ./cmd/sego --config ./sego.dev.json admin create
```

Frontend:

```bash
cd web
npm install
npm run dev
```

The development UI expects the API at `http://127.0.0.1:8088` unless configured otherwise.

## Security model

SEGO treats filenames, paths, uploaded content, share tokens, browser input, and external network traffic as untrusted. Core controls include:

- path normalization and storage-root confinement
- user namespace isolation
- server-side quota and authorization enforcement
- Argon2id password hashing
- opaque token generation and token hashing at rest where practical
- login throttling
- security response headers
- limited upload concurrency
- no automatic disk formatting or repartitioning

Before Internet exposure, read [`SECURITY.md`](SECURITY.md), [`docs/THREAT-MODEL.md`](docs/THREAT-MODEL.md), and [`docs/SECURITY-HARDENING.md`](docs/SECURITY-HARDENING.md).

## Repository layout

```text
cmd/                 Go entry points for server and CLI
internal/            Authentication, database, storage, sharing, API and config
web/                 React + TypeScript + Vite frontend
deploy/              Service/deployment templates
scripts/             Windows/Linux install and uninstall helpers
docs/                Technical and operational documentation
.github/workflows/    CI automation
```

## Documentation

Start with [`docs/README.md`](docs/README.md). Key documents include:

- [`ARCHITECTURE.md`](ARCHITECTURE.md)
- [`INSTALL.md`](INSTALL.md)
- [`SECURITY.md`](SECURITY.md)
- [`CONTRIBUTING.md`](CONTRIBUTING.md)
- [`UPGRADE.md`](UPGRADE.md)
- [`BACKUP-RESTORE.md`](BACKUP-RESTORE.md)
- [`TEST-REPORT.md`](TEST-REPORT.md)
- [`docs/API.md`](docs/API.md)
- [`docs/CLI.md`](docs/CLI.md)
- [`docs/CONFIGURATION.md`](docs/CONFIGURATION.md)
- [`docs/LOW-RESOURCE-MODE.md`](docs/LOW-RESOURCE-MODE.md)
- [`docs/TROUBLESHOOTING.md`](docs/TROUBLESHOOTING.md)

## Verification policy

SEGO does not claim tests, security properties, performance targets, signing, or platform support unless they have actually been executed and documented. The current workspace verification status is recorded in [`TEST-REPORT.md`](TEST-REPORT.md).

## Contributing

Contributions are welcome. Read [`CONTRIBUTING.md`](CONTRIBUTING.md) and the security reporting rules before submitting changes.

## License

SEGO Server is licensed under the MIT License. See [`LICENSE`](LICENSE).
