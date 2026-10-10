# Cloudflare Tunnel

SEGO can optionally expose only its web application through a Cloudflare Tunnel without opening an inbound router port.

## Beginner setup

First create a remotely-managed tunnel in the Cloudflare dashboard and configure its published application service as:

```text
http://localhost:8088
```

Then, on the SEGO server computer, run:

```text
sego tunnel setup
```

Paste the tunnel token when prompted. The prompt is hidden so the token is not intentionally added to command history.

## Boot-time reconnect

SEGO installs `cloudflared` using the OS service model:

- Linux: cloudflared systemd service
- Windows: cloudflared Windows service

Once configured, the tunnel starts/reconnects automatically after reboot.

Check it with:

```text
sego tunnel status
```

Remove the local service with:

```text
sego tunnel remove
```

## Security

A tunnel token is a secret. Anyone holding it can run a connector for that tunnel. If it may be exposed, rotate it in Cloudflare and reinstall the tunnel service with the new token.

Do not use the SEGO hostname to expose SQLite, SSH, database ports, RDP, or unrelated services unless you separately understand and secure those protocols. SEGO's intended public origin is only `http://localhost:8088` behind the tunnel/reverse proxy.
