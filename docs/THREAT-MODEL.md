# Threat Model

## Scope

This document describes security assumptions and major threat classes for the SEGO Server MVP architecture. It is not a substitute for penetration testing or deployment-specific risk assessment.

## Protected assets

- user file contents
- filenames and metadata
- account credentials and password hashes
- authenticated sessions
- share tokens and share policy
- SQLite database contents
- application configuration
- audit/diagnostic logs
- backup material
- update/release artifacts

## Trust boundaries

1. browser/client to HTTP API
2. authenticated user to another user's namespace
3. API server to configured storage roots
4. public share endpoint to protected file data
5. service process to operating-system resources
6. reverse proxy/tunnel to local SEGO listener
7. release/update source to installed executable

## Primary threats

### Path traversal and filesystem escape

An attacker may submit `../`, crafted separators, symlinks, or renamed paths to escape an authorized namespace.

**MVP controls:** path normalization, configured-root confinement, user namespace mapping, and symlink rejection in sensitive resolver paths.

### Cross-user IDOR

A user may attempt to access another user's object by guessing an ID or internal path.

**MVP controls:** authenticated user context and server-side namespace translation. Frontend visibility is never treated as authorization.

### Session/token theft

Leaked bearer or share tokens can grant access until expiration or revocation.

**MVP controls:** cryptographically random opaque tokens, hashing at rest where practical, expiration/revocation checks, and log avoidance.

### Credential attacks

Attackers may brute-force usernames/passwords.

**MVP controls:** Argon2id password hashing, generic credential errors, dummy verification work for unknown users, and login throttling.

### Quota race

Concurrent uploads may attempt to bypass quota checks.

**MVP controls:** transaction-oriented quota reservation and bounded upload concurrency.

### Malicious uploaded content

Files may contain scripts, macros, malformed media, decompression bombs, or payloads targeting preview libraries.

**Current position:** uploaded content is treated as untrusted; advanced preview isolation is not yet complete. Active content must not be executed merely for preview.

### Stored XSS

Filenames, Markdown, SVG, HTML, or metadata may contain browser-executable content.

**MVP controls:** restrictive security headers and no assumption that uploaded content is safe. Full preview sanitization/sandboxing remains roadmap work.

### Public tunnel misconfiguration

A tunnel/reverse proxy can unintentionally expose local-only services or weaken authentication.

**Controls/requirements:** keep application auth enabled, expose only the intended HTTP endpoint, never expose database/SSH through the share hostname, and validate proxy policy.

### Malicious updates

A compromised update source can deliver arbitrary code.

**Current position:** signed automatic update installation is not implemented. Unverified update execution is prohibited by project policy.

## Residual risks

Important incomplete controls include MFA/TOTP, complete folder ACL inheritance, full preview sandboxing, signed update/rollback verification, automated Cloudflare policy checks, comprehensive audit monitoring, native Windows service hardening, and restoration drills.

These items must remain visible in release notes until implemented and tested.
