# SEGO Server Test Report

**Workspace date:** 2026-10-10

This report records only checks that were actually executed in the current build workspace. It intentionally separates executed verification from CI that is configured but not yet observed here.

## Executed and passed

### Dependency-light Go package tests

```text
go test ./internal/config ./internal/namespace ./internal/storage
ok  github.com/nadeemmhdm/sego-server/internal/config
ok  github.com/nadeemmhdm/sego-server/internal/namespace
ok  github.com/nadeemmhdm/sego-server/internal/storage
```

Validated behaviors include:

- low-resource mode configuration
- standard mode configuration
- rejection of unknown runtime modes
- storage-root traversal rejection
- user namespace traversal confinement
- symbolic-link rejection on the supported Unix test path
- streaming write followed by directory listing
- namespace-to-public-path stripping

### Shell syntax

```text
bash -n scripts/install.sh scripts/uninstall.sh
```

Result: passed.

## Not executable in this sandbox

### Complete Go dependency-backed test suite

The full repository depends on external modules including Chi, `x/crypto`, `x/term`, and `modernc.org/sqlite`. Earlier attempts in this sandbox could not fetch missing modules because outbound module download access was unavailable. Therefore this report does not claim `go test ./...` passed here.

### PowerShell parser/runtime validation

PowerShell (`pwsh`) is not installed in the current Linux workspace, so the Windows installer was reviewed and prepared for Windows CI but not executed here.

### Frontend dependency build

The web build requires npm dependency installation. This report does not claim a production frontend build was completed in this workspace unless a later CI result explicitly records it.

## Still requiring broader validation

- Windows binary build and native service behavior
- Ubuntu/Debian clean-machine install/reinstall
- full SQLite/auth/share integration tests
- concurrent quota race integration tests
- interrupted upload cleanup
- disk-full behavior
- Cloudflare Tunnel lifecycle
- update signature/rollback path
- backup restoration drills
- real-hardware 4 GB memory/CPU benchmarks
- arm64 validation

## CI policy

The repository includes GitHub Actions intended to run the complete dependency-backed test/build matrix after the repository is published. A configured workflow is not itself evidence that a test passed; release notes should reference observed results.
