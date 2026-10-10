# Updates

SEGO uses GitHub Releases as its update channel.

## Commands

```text
sego version
sego update check
sego update auto
```

## Automatic schedule

Normal one-command installation enables background update checks:

- Linux: `sego-update.timer`, after boot and approximately daily
- Windows: `SEGO Update` scheduled task, at startup and daily

If no published GitHub Release exists, no update is installed.

## Update package requirements

A release is installable only when it contains:

- a platform archive named `sego-server_<version>_<os>_<arch>.tar.gz` or `.zip`
- `SHA256SUMS` containing the archive hash

The updater downloads both before stopping SEGO and verifies SHA-256 before extraction.

## Rollback behavior

The updater keeps the current binaries and dashboard in a temporary local backup while installing a release. After restarting, it calls the local health endpoint. If the new server does not become healthy, the updater attempts to restore the previous files and restart the old server.

## Release workflow

Pushing a `v*` Git tag triggers the repository release workflow. It builds platform packages, generates `SHA256SUMS`, and publishes GitHub Release assets.

## Current trust boundary

The current MVP relies on HTTPS/GitHub release delivery plus SHA-256 integrity checks. A separately signed release manifest / pinned project signing key is still roadmap work and should be completed before high-assurance production auto-update claims.
