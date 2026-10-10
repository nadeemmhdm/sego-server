# Installation Guide

## Source-install prerequisites

The current installers build SEGO from source and therefore require:

- Go 1.23+
- a supported source checkout
- administrator/root privileges for system-wide installation

The frontend build is separate in the current MVP.

## Linux

Initially targeted distributions are Ubuntu Server and Debian with systemd.

Low-resource profile:

```bash
sudo ./scripts/install.sh --profile low-resource --yes
```

Standard profile:

```bash
sudo ./scripts/install.sh --profile standard --yes
```

Default locations:

```text
/opt/sego/bin            binaries
/etc/sego/sego.json      configuration
/var/lib/sego            database and primary storage
/etc/systemd/system      service unit
```

Create the first administrator after installation:

```bash
sudo /opt/sego/bin/sego --config /etc/sego/sego.json admin create
```

Check health:

```bash
sudo /opt/sego/bin/sego --config /etc/sego/sego.json doctor
systemctl status sego.service
```

## Windows

Run PowerShell as Administrator from the repository root:

```powershell
powershell -ExecutionPolicy Bypass -File .\scripts\install.ps1 -Profile low-resource -Yes
```

Default locations:

```text
%ProgramFiles%\SEGO\bin
%ProgramData%\SEGO\sego.json
%ProgramData%\SEGO\storage
```

Create the administrator:

```powershell
& "$env:ProgramFiles\SEGO\bin\sego.exe" --config "$env:ProgramData\SEGO\sego.json" admin create
```

Start the MVP server manually:

```powershell
& "$env:ProgramFiles\SEGO\bin\sego-server.exe" --config "$env:ProgramData\SEGO\sego.json"
```

Native Windows Service registration is intentionally documented as incomplete until the service-control integration is implemented and tested. The installer does not fake service success.

## Uninstall

Linux, preserving data by default:

```bash
sudo ./scripts/uninstall.sh
```

Linux, including data/configuration:

```bash
sudo ./scripts/uninstall.sh --purge-data
```

Windows, preserving data:

```powershell
.\scripts\uninstall.ps1
```

Windows, including data/configuration:

```powershell
.\scripts\uninstall.ps1 -PurgeData
```

Back up important data before any destructive uninstall option.
