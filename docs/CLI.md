# CLI Reference

The CLI executable is named `sego`.

## Global option

```text
--config PATH   Path to the SEGO JSON configuration file
```

## Implemented MVP commands

### Status

```bash
sego --config sego.json status
```

Shows the listen address, database path, storage-root count, and active resource mode.

### Doctor

```bash
sego --config sego.json doctor
```

Checks configuration and storage-root accessibility.

### Storage list

```bash
sego --config sego.json storage list
```

Lists configured storage roots.

### Create administrator

```bash
sego --config sego.json admin create
```

Interactively creates a Super Admin. Password input is not designed to be passed through command-line arguments.

### Show mode

```bash
sego --config sego.json mode show
```

### Set low-resource mode

```bash
sego --config sego.json mode set low-resource
```

### Set standard mode

```bash
sego --config sego.json mode set standard
```

## Planned commands

The project specification also calls for tunnel setup/status, update check/apply, backup/restore, logs, and full uninstall management. These remain roadmap items until implemented and tested.
