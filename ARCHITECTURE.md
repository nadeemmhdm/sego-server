# SEGO Server Architecture

## Purpose

SEGO Server uses a deliberately small, cross-platform architecture suitable for self-hosted storage on modest hardware. The storage and authorization layers remain independent from the browser UI so security checks cannot be bypassed by directly calling APIs.

## Components

### API server

The Go API server owns authentication, authorization, storage access, quota checks, sharing, diagnostics, and HTTP security controls.

### CLI

The `sego` CLI provides local administration and diagnostics. Current MVP commands cover status, health checks, storage listing, first-admin creation, and runtime mode management.

### Database

SQLite stores application metadata such as users, password hashes, sessions, file index data, quota reservations, and share records. File contents are not stored in SQLite.

### Storage manager

File contents remain on explicitly configured native filesystem roots. Storage access is resolved through a confinement layer that prevents directory escape and unsafe cross-user access.

### User namespaces

Authenticated file operations are mapped into a per-user internal namespace. Public API paths are not trusted as literal filesystem paths.

### Web application

The React/TypeScript frontend communicates with the API. It is not a security boundary; all permissions and quotas are enforced by the backend.

## Data flow

```text
Browser / CLI
     |
     v
HTTP API / local CLI
     |
     +--> Authentication and authorization
     |
     +--> SQLite metadata / sessions / shares
     |
     +--> Storage manager --> configured filesystem roots
```

## Resource strategy

SEGO uses streaming transfers, bounded upload concurrency, small database pools, pagination-oriented APIs, and a low-resource mode. Large file contents should never be loaded fully into memory during normal transfer operations.

## Trust boundaries

Primary trust boundaries include the browser-to-API boundary, authenticated user-to-user separation, application-to-filesystem access, public share endpoints, update packages, and remote tunnel exposure.

See [`docs/THREAT-MODEL.md`](docs/THREAT-MODEL.md) for security-specific analysis.
