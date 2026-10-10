# Troubleshooting

## `sego doctor` fails on a storage root

Verify the configured path exists and that the SEGO service identity can read/write it. Do not solve permission errors by granting broad world-writable access.

## Linux service does not start

```bash
sudo systemctl status sego.service
sudo journalctl -u sego.service -n 100 --no-pager
```

Confirm the config path and data paths match the generated service unit.

## Port 8088 is already in use

Change `listen_addr` to another local port and restart the service. Confirm the replacement port is not unintentionally public.

## Upload rejected

Possible causes include maximum request size, quota limit, disk-full state, storage-root availability, or upload concurrency limit.

## Login rejected

Check the username/password, account state, and server clock. Repeated failures may be rate-limited.

## Windows server does not start automatically

Automatic native Windows Service registration is not implemented in the MVP source installer. Start the server manually until the SCM integration is completed and tested.
