# Cloudflare Tunnel

SEGO can optionally expose only its web application through a Cloudflare Tunnel without opening an inbound router port.

## Beginner setup

First create a remotely-managed tunnel in the Cloudflare dashboard and configure its published application service as:

```text
http://localhost:8088
```

Then, on the SEGO server computer, run:

```text
sego remote
```

`sego remote` is the beginner alias for `sego tunnel setup`. Paste the tunnel token when prompted. The prompt is hidden so the token is not intentionally added to command history.

## Automatic boot reconnect

Tunnel setup is needed only once unless the token is rotated or the service is removed.

After setup:

- Linux installs/enables the `cloudflared` systemd service, orders it after networking/SEGO, and configures automatic restart after failures.
- Windows configures the `cloudflared` Windows service for automatic boot startup and restart-on-failure recovery.
- If the laptop/server is powered off and later turned on, SEGO starts in the background and the configured tunnel reconnects automatically when networking is available.
- If no tunnel has ever been configured, SEGO does not invent a tunnel or expose the server automatically.

Check it with:

```text
sego tunnel status
```

Remove the local tunnel service with:

```text
sego tunnel remove
```

## Security

A tunnel token is a secret. Anyone holding it can run a connector for that tunnel. If it may be exposed, rotate it in Cloudflare and reinstall the tunnel service with the new token.

Do not use the SEGO hostname to expose SQLite, SSH, database ports, RDP, or unrelated services unless you separately understand and secure those protocols. SEGO's intended public origin is only `http://localhost:8088` behind the tunnel/reverse proxy.
