# Getting Started

## 1. Choose a mode

Use `low-resource` for older or 4 GB-class systems. Use `standard` when the host has more CPU/RAM headroom.

## 2. Install

Linux:

```bash
sudo ./scripts/install.sh --profile low-resource --yes
```

Windows (elevated PowerShell):

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Profile low-resource -Yes
```

## 3. Create the administrator

Linux:

```bash
sudo /opt/sego/bin/sego --config /etc/sego/sego.json admin create
```

Windows:

```powershell
& "$env:ProgramFiles\SEGO\bin\sego.exe" --config "$env:ProgramData\SEGO\sego.json" admin create
```

## 4. Verify health

```bash
sego --config /path/to/sego.json doctor
```

## 5. Keep it local first

Validate authentication, file listing, upload/download, quota behavior, and backups before adding any public tunnel or reverse proxy.

## 6. Read the security docs

Review `SECURITY.md`, `docs/THREAT-MODEL.md`, and `docs/SECURITY-HARDENING.md` before exposing SEGO outside a trusted LAN.
