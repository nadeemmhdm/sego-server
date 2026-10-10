# Configuration Reference

SEGO uses a JSON configuration file.

## Core fields

| Field | Purpose |
|---|---|
| `listen_addr` | API listen address |
| `database_path` | SQLite database location |
| `session_hours` | authenticated session lifetime |
| `max_upload_bytes` | per-request upload size ceiling |
| `upload_concurrency` | maximum concurrent upload handlers |
| `low_resource_mode` | enables conservative runtime behavior |
| `storage_roots` | explicitly configured native filesystem roots |

## Safe defaults

The default listener is `127.0.0.1:8088`, avoiding direct public exposure.

The default Linux data path is `/var/lib/sego`; Windows uses `%ProgramData%\SEGO` when available.

## Runtime mode commands

```bash
sego --config /path/to/sego.json mode show
sego --config /path/to/sego.json mode set low-resource
sego --config /path/to/sego.json mode set standard
```

Configuration updates are written through a temporary file and rename operation to reduce the chance of leaving a partially written configuration.

## Storage roots

Every root should point to an existing filesystem location controlled by the operator. Do not use a sensitive OS directory as a storage root. The server applies confinement checks, but host filesystem permissions should still follow least privilege.
