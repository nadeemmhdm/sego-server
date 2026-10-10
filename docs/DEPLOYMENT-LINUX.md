# Linux Deployment

## Initial target

SEGO initially targets Ubuntu Server and Debian using systemd.

## Service account

The installer creates a dedicated `sego` system account with a non-login shell. User data is stored separately from the application binaries.

## systemd hardening

The supplied unit uses:

- `NoNewPrivileges=true`
- `PrivateTmp=true`
- `ProtectSystem=strict`
- `ProtectHome=true`
- an explicit writable path set
- restrictive umask

Review writable paths whenever storage roots are customized.

## Health checks

```bash
sudo /opt/sego/bin/sego --config /etc/sego/sego.json doctor
sudo systemctl status sego.service
sudo journalctl -u sego.service -n 100 --no-pager
```

## Network exposure

Keep the API bound to loopback unless the deployment design explicitly requires otherwise. Terminate TLS at a trusted reverse proxy/tunnel and keep application authentication enabled.
