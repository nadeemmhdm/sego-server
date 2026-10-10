# Contributing to SEGO Server

Thank you for contributing to SEGO Server. Changes should preserve the project's priorities: data safety, explicit authorization, low resource consumption, maintainability, and honest verification.

## Development setup

1. Install Go 1.23+ and Node.js 22+.
2. Fork or create a working branch.
3. Run `go test ./...` before submitting backend changes.
4. Run `npm install` and `npm run build` in `web/` before submitting frontend changes.
5. Add or update tests when behavior changes.
6. Update documentation for user-visible or operational changes.

## Change requirements

Security-sensitive changes must consider path traversal, IDOR/cross-user access, direct URL bypass, malicious filenames, token leakage, quota races, stored XSS, request-size limits, and rollback/recovery behavior.

Do not:

- weaken backend authorization in favor of frontend-only checks
- add default credentials
- log passwords or full share tokens
- introduce unbounded worker pools
- read entire large uploads into memory
- silently format or repartition storage
- claim a benchmark or test passed unless it was executed

## Commit style

Use concise, scoped messages where practical, for example:

```text
feat(storage): add resumable upload session metadata
fix(auth): revoke expired sessions during validation
docs(security): document reverse proxy hardening
test(sharing): cover revoked link access
```

## Pull requests

Describe the problem, implementation, tests performed, security implications, resource implications, and any migration or rollback requirements.

Security vulnerabilities should follow [`SECURITY.md`](SECURITY.md), not the public issue tracker.
