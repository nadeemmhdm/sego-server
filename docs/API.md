# API Reference — MVP

Base development address: `http://127.0.0.1:8088`

Except for health, login, and authorized public share access, protected endpoints require:

```http
Authorization: Bearer <session-token>
```

## Public endpoints

### `GET /api/health`

Returns basic service health.

### `POST /api/auth/login`

Authenticates a user and returns an opaque session token.

### `GET /api/share/{token}/download`

Downloads an authorized shared file when the share is valid, not revoked, not expired, and within configured access policy.

## Authenticated endpoints

### `POST /api/auth/logout`

Revokes the current session.

### `GET /api/me`

Returns the current authenticated user.

### `GET /api/files`

Lists an authorized directory within the current user's namespace.

### `POST /api/files/mkdir`

Creates a directory.

### `POST /api/files/upload`

Streams an upload subject to request-size, quota, concurrency, namespace, and storage-root checks.

### `GET /api/files/download`

Streams an authorized file download.

### `POST /api/files/rename`

Renames/moves an authorized item within permitted storage scope.

### `DELETE /api/files`

Moves an authorized item to trash according to the current MVP behavior.

### `POST /api/shares`

Creates a share record and returns the plaintext token once to the caller while persisting a token hash where practical.

### `DELETE /api/shares/{id}`

Revokes a share owned by the authorized user/administrator according to backend policy.

### `GET /api/diagnostics`

Returns resource/service diagnostics available to the application.

## Security requirements

Clients must never assume a path is authorized because it appeared in the UI. The server remains responsible for authorization, namespace translation, quota enforcement, share validation, and path confinement.
