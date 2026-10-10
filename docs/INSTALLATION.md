# Installation Reference

For beginner installation, use the one-command instructions in [`GETTING-STARTED.md`](GETTING-STARTED.md).

## Direct GitHub bootstrap

Windows:

```powershell
irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
```

Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
```

The bootstrap resolves the current `main` commit, downloads that source snapshot to a temporary directory, and runs the platform-local installer from the downloaded snapshot. The Linux bootstrap should be started as a normal user rather than with sudo.

## Automatic profile selection

`auto` is the default. Systems below approximately 6 GiB RAM use `low-resource`; larger systems use `standard` unless the operator explicitly chooses a profile.

## Linux layout

- binaries: `/opt/sego/bin`
- simple CLI symlink: `/usr/local/bin/sego`
- config: `/etc/sego/sego.json`
- application data: `/var/lib/sego`
- service: `sego.service`

## Windows layout

- binaries: `%ProgramFiles%\SEGO\bin`
- config/data: `%ProgramData%\SEGO`
- startup: `SEGO Server` Scheduled Task in the current MVP
- binary directory is added to machine PATH

## Source checkout installation

Developers/advanced operators can run `scripts/install.sh` or `scripts/install.ps1` directly from a checkout. See each script's help/options.
