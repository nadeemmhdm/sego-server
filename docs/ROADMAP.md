# SEGO Server Roadmap

Roadmap status reflects implementation maturity rather than intended scope.

## Phase 1 — Foundation

**Status: implemented foundation**

- architecture and repository structure
- initial threat model
- Windows/Linux support strategy
- low-resource design principles

## Phase 2 — Installation, CLI and admin bootstrap

**Status: partial / MVP**

- Go CLI
- administrator creation
- Linux systemd source installer
- Windows source installer
- runtime mode switching

Remaining: native Windows Service integration, signed release installer, richer non-interactive bootstrap.

## Phase 3 — Storage and quotas

**Status: MVP implemented**

- configured storage roots
- list/upload/download/mkdir/rename/trash
- per-user namespace
- streaming file operations
- quota reservation logic

Remaining: richer metadata indexing, multiple-root operational workflows, removable-drive handling and storage reports.

## Phase 4 — Sharing and permissions

**Status: partial**

- opaque hashed file-share tokens
- expiration and revocation foundation

Remaining: registered-user/group sharing, complete folder sharing, inheritance rules, password-protected links, recipient verification and advanced download policy.

## Phase 5 — Preview, sorting and search

**Status: pending expansion**

- safe preview pipeline
- thumbnails
- richer server-side sort/filter/search
- optional content indexing

## Phase 6 — Remote access

**Status: pending**

- guided Cloudflare Tunnel lifecycle
- health/diagnostics
- optional Cloudflare Access integration
- configuration validation

## Phase 7 — UI and accessibility

**Status: partial**

- React/TypeScript/Vite dashboard foundation
- responsive layout foundation

Remaining: complete file manager interactions, polished motion, low-resource visual reductions, advanced accessibility verification.

## Phase 8 — Updates, backup and recovery

**Status: pending**

- signed update verification
- migration/rollback orchestration
- scheduled backups
- encrypted archive option
- restore validation workflows

## Phase 9 — Release hardening

**Status: pending**

- expanded security regression suite
- Windows/Linux installation matrix
- real-hardware 4 GB benchmarks
- release signing
- published measured performance results
