# Installation Details

## Installation philosophy

SEGO installation must be idempotent where practical, must not format or repartition disks, and must not run unverified remote scripts with elevated privileges.

The current MVP source installers build from the local checkout. They do not download a claimed release binary that has not been signed and published.

## Linux flow

The Linux installer:

1. validates Linux and architecture
2. identifies the distribution when possible
3. verifies Go and systemd are available
4. shows the install plan
5. builds `sego-server` and `sego`
6. creates the restricted `sego` service account
7. creates application directories
8. initializes configuration
9. applies the selected resource mode
10. installs/enables the systemd service
11. runs diagnostics and service health checks

## Windows flow

The Windows installer:

1. verifies elevated PowerShell
2. verifies Go is available
3. shows the install plan
4. builds `sego-server.exe` and `sego.exe`
5. initializes configuration
6. applies the resource mode
7. runs diagnostics

Native Windows Service installation is still pending implementation and validation.

## Profiles

`low-resource`:

- low-resource mode enabled
- upload concurrency reduced to 1
- recommended default on older hardware

`standard`:

- low-resource mode disabled
- upload concurrency raised to at least 3
- intended for systems with more resource headroom
