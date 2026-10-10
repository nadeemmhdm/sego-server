# SEGO Server Installation

For most users, installation is one command. You do not need to download the source code manually.

## Windows 10 / 11

1. Open PowerShell.
2. Paste this command:

```powershell
irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
```

3. Choose **Yes** when Windows asks for Administrator permission.
4. Create your SEGO administrator username and password when prompted.
5. The installer starts SEGO automatically and opens the dashboard.

Afterwards:

```text
sego open
sego status
sego start
sego stop
```

If `sego` is not recognized in the same PowerShell window immediately after the first install, open a new PowerShell window once so Windows reloads the updated PATH.

## Ubuntu / Debian

Run as your normal user:

```bash
curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
```

Do **not** put `sudo` in front of the command. SEGO asks for sudo only for package/service installation.

The installer detects lower-memory machines automatically and selects Low Resource Mode when appropriate.

## What is automatic

The guided installer handles:

- GitHub source download;
- supported dependency installation;
- verified Go toolchain installation on Linux when the installed Go version is too old;
- frontend/backend build;
- config/data directories;
- automatic startup;
- first administrator setup;
- health checks;
- PATH/launcher setup.

## Interrupted first setup

If installation finished but the administrator prompt was skipped:

```text
sego setup
```

On Linux, if the protected database requires elevation:

```bash
sudo sego setup
```

On Windows, open an Administrator PowerShell only for this recovery step.

## Custom/advanced installation

The repository-local installers still support advanced paths and profiles:

Linux:

```bash
./scripts/install.sh --profile low-resource
```

Windows:

```powershell
.\scripts\install.ps1 -Profile low-resource
```

Advanced options are intended for administrators and development environments, not normal first-time setup.
