# Windows Deployment

## Initial target

SEGO targets tested Windows 10/11 configurations.

## Source install

Run PowerShell as Administrator:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Profile low-resource -Yes
```

## Current service status

The MVP source installer installs binaries and configuration but does not claim native Windows Service registration. Native SCM integration, least-privileged service identity, recovery behavior, and signing must be implemented and tested before that feature is marked complete.

## Manual start

```powershell
& "$env:ProgramFiles\SEGO\bin\sego-server.exe" --config "$env:ProgramData\SEGO\sego.json"
```

## Permissions

Protect `%ProgramData%\SEGO` so ordinary users cannot read database/session/share metadata or modify configuration. Production installation should apply explicit ACLs as part of the future native installer.
