# Updates

SEGO uses published GitHub Releases as its normal automatic update channel.

## Beginner commands

```text
sego version
sego update check
sego update
```

`sego update` installs the newest eligible published release. The longer compatibility command `sego update auto` still works.

## Automatic schedule

Normal one-command installation enables background update checks without requiring the user to open SEGO:

- Linux: `sego-update.timer` starts a few minutes after boot, then checks approximately every 6 hours with a small randomized delay.
- Windows: `SEGO Update` runs at startup and daily, uses `StartWhenAvailable`, and retries transient failures.

The updater calls GitHub's latest-release endpoint. Draft and prerelease builds are not selected by that endpoint. If no newer published GitHub Release exists, nothing is installed.

## Update package requirements

A release is installable only when it contains:

- a platform archive named `sego-server_<version>_<os>_<arch>.tar.gz` or `.zip`
- `SHA256SUMS` containing the archive hash

The updater downloads both before stopping SEGO and verifies the package SHA-256 before extraction. After replacement it restarts SEGO and verifies the local health endpoint.

## Boot and network behavior

SEGO itself starts as a background boot service/task. Update checks do not require a browser session or user login. If networking is temporarily unavailable, Linux checks again on the next timer activation and Windows is configured to retry/start when available.

## Rollback behavior

The updater keeps the current binaries and dashboard in a temporary local backup while installing a release. After restarting, it calls the local health endpoint. If the new server does not become healthy, the updater attempts to restore the previous files and restart the old server.

## Release workflow

Pushing a `v*` Git tag triggers the repository release workflow. It builds platform packages, embeds the tag version into the binaries, generates `SHA256SUMS`, and publishes GitHub Release assets.

## Current trust boundary

The current MVP relies on HTTPS/GitHub release delivery plus SHA-256 integrity checks. A separately signed release manifest / pinned project signing key is still roadmap work and should be completed before high-assurance production auto-update claims.
