# API Reference — MVP

Default local base address: `http://127.0.0.1:8088`

Protected endpoints require:

```http
Authorization: Bearer <session-token>
```

During Emergency SOS lockdown, login, authenticated endpoints, public shares, and preview-token access are blocked. Health remains available for local service diagnostics.

## Public endpoints

### `GET /api/health`

Basic service health.

### `POST /api/auth/login`

Authenticates a user and returns an opaque session token. Login is rate-limited and disabled during emergency lockdown.

### `GET /api/share/{token}/download`

Downloads a valid authorized share. The backend validates revocation, expiry and configured limits.

### `GET /api/preview/{token}`

Streams an authorized temporary preview. Media uses HTTP Range support through `http.ServeContent`. Tokens are scoped to one stored file and expire automatically.

## Authenticated endpoints

### `POST /api/auth/logout`

Revokes the current session.

### `GET /api/me`

Returns the current user.

### `GET /api/files`

Lists an authorized directory in the current user's namespace.

### `POST /api/files/mkdir`

Creates a directory.

### `POST /api/files/upload`

Streams an upload subject to request size, quota, concurrency and path-confinement checks.

### `GET /api/files/download`

Streams an authorized download.

### `POST /api/files/preview-token`

Creates a temporary scoped preview URL for an explicitly supported safe preview type. Active HTML/SVG is not admitted as inline preview content.

### `POST /api/files/rename`

Renames/moves an item within the user's allowed storage namespace.

### `DELETE /api/files`

Moves an authorized item to SEGO trash according to current MVP behavior.

### `POST /api/shares`

Creates a file share. Plaintext token is returned to the caller while a hash is persisted where practical.

### `DELETE /api/shares/{id}`

Revokes an authorized share.

### `GET /api/diagnostics`

Admin-only resource/service diagnostics.

## Admin security endpoints

### `GET /api/admin/security/status`

Returns the global Emergency SOS state. Requires `admin` or `super_admin`.

### `POST /api/admin/emergency/activate`

Activates global emergency lockdown. This operation:

- revokes every session;
- revokes/expires existing share links;
- deletes preview tokens;
- blocks further login/authenticated access;
- writes an audit event.

There is intentionally no remote unlock endpoint. Unlock is local CLI + administrator-password recovery.

### `GET /api/admin/audit?limit=100&offset=0`

Returns recent audit events, including action, timestamp, actor details, IP, user agent, target, outcome and structured details. Requires `admin` or `super_admin`.

## Security requirements

Clients must never infer authorization from UI visibility. Namespace translation, permission checks, quota enforcement, share validation, preview validation, emergency state and filesystem confinement remain server-side responsibilities.
