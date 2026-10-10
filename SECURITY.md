# Security Policy

Security is a primary design constraint for SEGO Server.

## Supported security scope

The project is currently an MVP foundation. Security reports are welcome for the latest default branch. Release-specific support windows will be documented after signed versioned releases begin.

## Reporting a vulnerability

Do **not** publish exploit details, credentials, share tokens, personal data, or active-server information in a public issue.

Use GitHub's private security reporting feature when it is enabled for the repository. If private reporting is unavailable, contact the repository owner through a private channel and include only the minimum information needed to reproduce the issue safely.

A useful report includes:

- affected component and version/commit
- deployment platform
- prerequisites
- reproduction steps
- expected versus actual result
- security impact
- sanitized logs or proof of concept
- suggested mitigation, if known

## Security principles

SEGO follows these baseline rules:

- no plaintext password storage
- Argon2id password hashing
- opaque random session/share tokens
- token hashing at rest where practical
- authorization enforced on the server
- storage-root confinement
- user namespace isolation
- no raw filesystem paths in public URLs
- bounded request/upload processing
- no automatic disk formatting or repartitioning
- no logging of passwords or full share tokens
- no claim that a link is secure merely because it is difficult to guess
- no unverified update package execution
- scoped temporary preview tokens instead of bearer tokens in media URLs
- Emergency SOS global session/link/preview revocation
- administrator audit log with actor, IP, outcome and timestamp

## Deployment responsibility

Operators are responsible for TLS termination, host patching, filesystem permissions, secure backup storage, Cloudflare or reverse-proxy policy, firewall rules, account lifecycle, and physical security.

Do not expose database, SSH, or local administrative ports through a public file-sharing hostname.

## Known incomplete controls

MFA/TOTP, complete folder ACL inheritance, signed update/rollback automation, native Windows service hardening, centralized alerting/SIEM export, and complete security-regression coverage remain roadmap items. Do not assume they are already present.

See [`docs/SECURITY-HARDENING.md`](docs/SECURITY-HARDENING.md) and [`docs/THREAT-MODEL.md`](docs/THREAT-MODEL.md).
