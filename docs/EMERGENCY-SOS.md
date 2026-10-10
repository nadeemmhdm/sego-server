# Emergency SOS

Emergency SOS is a containment control for suspected account compromise, leaked links, unexpected access, or other active security incidents.

## Activation

An administrator can activate SOS from **Dashboard → Security → Emergency SOS**.

Activation is intentionally immediate. It performs one transaction that:

1. sets global lockdown active;
2. revokes all existing login sessions, including the activating admin;
3. revokes and expires all active share links;
4. deletes temporary preview tokens;
5. records the action, actor, IP, user agent, reason and timestamp in the audit log.

After activation, remote login and protected API access return a locked response. Public shares and preview URLs are also blocked.

## Recovery

Emergency recovery is deliberately not exposed through a remote API.

On the physical/server host run:

```text
sego emergency unlock
```

The CLI asks for an administrator username and password. On protected Linux/Windows deployments, host elevation may be required to access the database.

## Post-recovery behavior

- Users must sign in again.
- Old sessions are not restored.
- Old share links stay revoked.
- Old preview tokens stay invalid.
- Owners must create new share links when required.

This behavior avoids silently reopening access that was intentionally terminated during an incident.
