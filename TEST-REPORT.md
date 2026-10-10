# SEGO Server Verification Report

Date: 2026-10-10

This report distinguishes checks actually executed in the development workspace from checks delegated to GitHub CI or still pending.

## Executed successfully in the workspace

- `bash -n install.sh scripts/install.sh scripts/bootstrap.sh`
- `go test ./internal/config`
- `go test ./internal/namespace`
- `go test ./internal/storage`
- package compilation checks for `internal/securitylock` and `internal/audit` through `go test`
- Go formatting for modified backend/CLI packages

The storage/namespace tests include root-confinement and traversal-related cases already present in the repository.

## Not fully executed in this workspace

### Full Go dependency build

`go test ./...` could not complete because the execution environment could not obtain missing external module checksum/dependency data. It reported missing `go.sum` entries for external modules. The repository CI is configured to run the full Go test/vet/build matrix with network-enabled GitHub runners.

### Frontend dependency build

`npm install` could not complete in the local execution sandbox because external package retrieval timed out. GitHub CI is configured to install pinned direct frontend dependencies and run the TypeScript/Vite build.

### PowerShell parser

PowerShell (`pwsh`) is not installed in the local Linux workspace. GitHub CI parses the Windows installer/bootstrap/uninstaller scripts using PowerShell.

## Toolchain correction

The repository dependencies require Go 1.24 or newer. The project `go` directive, CI toolchain and installer minimum have been aligned accordingly. The Linux guided installer can install a checksum-verified Go toolchain when the system Go is too old.

## Still required before production claims

- successful GitHub CI on the final commit
- Windows 10/11 clean-machine installation test
- Ubuntu/Debian clean-machine installation test
- browser preview/Range playback matrix
- Emergency SOS end-to-end recovery test
- concurrent quota and transfer stress tests
- full security regression suite
- measured 4 GB hardware resource benchmark
- signed release/update/rollback validation

No claim is made that those checks passed until results are captured.
