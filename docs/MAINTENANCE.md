# Maintenance Guide

## Routine checks

Operators should periodically verify:

- `sego doctor` succeeds
- expected storage roots are mounted and writable
- free disk capacity remains above operational thresholds
- backups are completing and restoration remains possible
- operating-system security updates are current
- administrator and user accounts remain necessary
- public shares are still required and appropriately scoped

## Logs

On Linux, the systemd journal is the primary service log source:

```bash
sudo journalctl -u sego.service --since today
```

Avoid copying logs to public issues without reviewing them for filenames, usernames, hostnames, IP addresses, or other sensitive metadata.

## Database

Do not manipulate the SQLite database directly while SEGO is actively writing unless the operation is explicitly documented as safe. Back up before schema-level maintenance.

## Storage

Never move or replace a configured storage root without verifying ownership, mount state, available capacity, and the effect on indexed metadata.
