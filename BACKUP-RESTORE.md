# Backup and Restore

The automated scheduled backup engine remains a roadmap item. Operators should maintain independent backups of both metadata and file content.

## Minimum backup set

Back up:

- SEGO configuration
- SQLite database
- configured storage roots
- deployment/service configuration

## Consistency

For a simple offline backup, stop SEGO before copying the SQLite database and storage metadata. Online database-consistent backup orchestration will be added later.

## Restore validation

A backup is not considered reliable until restoration has been tested. Validate:

- database opens successfully
- administrator login works
- expected storage roots are present
- representative files can be downloaded
- share state is correct
- permissions have not broadened unexpectedly

Keep at least one backup on a different physical device from the primary data.
