# Upgrade Policy

The signed automatic update/rollback manager is a roadmap item. Until it is implemented and verified, upgrades should be treated as controlled source deployments.

Before upgrading:

1. stop write activity
2. back up the SQLite database and configuration
3. back up critical user files or verify an independent backup
4. record the currently deployed commit/version
5. review schema or configuration changes
6. build and test the new version
7. restart SEGO and run `sego doctor`
8. verify login, listing, upload, download, and sharing

Do not perform destructive database or filesystem migration without a tested rollback path.

See [`BACKUP-RESTORE.md`](BACKUP-RESTORE.md) and [`docs/RELEASES.md`](docs/RELEASES.md).
