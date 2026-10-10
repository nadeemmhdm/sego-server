# Development Guide

## Requirements

- Go 1.23+
- Node.js 22+
- npm
- Git

## Backend

```bash
go test ./...
go run ./cmd/server --config ./sego.dev.json
```

## CLI

```bash
go run ./cmd/sego --config ./sego.dev.json status
go run ./cmd/sego --config ./sego.dev.json mode show
```

## Frontend

```bash
cd web
npm install
npm run dev
```

## Engineering rules

- backend is the authorization boundary
- preserve streaming for large file operations
- use bounded concurrency
- add tests for security-sensitive path and permission changes
- avoid introducing always-running helper services without clear justification
- document migrations and rollback implications
