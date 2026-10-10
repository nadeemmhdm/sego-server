# Audit Logs

SEGO stores security-relevant and file-operation events in SQLite for administrator review.

## Dashboard

Open **Audit Logs** from the admin sidebar. The list shows:

- action;
- user/role;
- IP address;
- date and time;
- outcome.

Select a row to view:

- event ID;
- exact action name;
- UTC timestamp and local rendering;
- actor user ID/name/role;
- client IP;
- user agent;
- target/resource;
- outcome;
- structured JSON details.

## Current audited events

The MVP records events including login success/failure, logout, file/folder mutations, share creation/revocation/download, preview issuance/opening, and Emergency SOS activation/recovery.

## Logging safety

SEGO must not write plaintext passwords or complete secret/share/preview tokens to audit records. HTTP request logging uses route patterns for tokenized endpoints so opaque token values are not intentionally copied into application logs.
