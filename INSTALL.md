# SEGO Server Installation

SEGO is designed so a first-time user can install it without cloning the repository or manually downloading source files.

## Windows 10 / 11

1. Open **PowerShell**.
2. Paste this one command:

```powershell
irm https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.ps1 | iex
```

3. Approve the Administrator/UAC prompt.
4. Enter the first admin username and password when asked.
5. Wait for **Installation complete**. The dashboard starts automatically.

After that, use:

```text
sego          Open dashboard
sego status   Check server
sego update   Update SEGO
sego help     Help
```

If Windows does not recognize `sego` in the same PowerShell window immediately after installation, close that window and open PowerShell once more so PATH reloads.

## Ubuntu / Debian

Open Terminal as your normal user and paste:

```bash
curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
```

Do **not** prefix it with `sudo`. The installer asks for sudo only for required OS changes.

Then use:

```text
sego          Open dashboard
sego status   Check server
sego update   Update SEGO
sego help     Help
```

## What is automatic

The installer handles GitHub download, supported prerequisite installation, verified Go toolchain fallback on Linux when required, performance-mode selection, Go dependency resolution, frontend/backend build, protected directories, startup service/task creation, automatic update checks, first-admin setup and health validation.

It does **not** format or repartition disks.

## If first-admin setup was interrupted

Run:

```text
sego setup
```

Linux may require `sudo sego setup` if the protected database cannot be opened by the current user.

## Advanced installation

Repository-local installers still expose profile/path options for administrators and development environments. Normal users do not need them.
