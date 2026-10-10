# SEGO CLI

The CLI is designed so normal users do not need to know configuration paths.

## Everyday commands

```text
sego open       Start SEGO if necessary and open the dashboard
sego status     Show local server/service status
sego start      Start the background server
sego stop       Stop the background server
sego restart    Restart the background server
sego help       Show help
```

On Windows, service-control commands can display the normal UAC prompt when elevation is required. On Linux, sudo can request the user's password when service control requires it.

## Setup / diagnostics

```text
sego setup
sego doctor
sego storage list
```

`sego setup` creates the first administrator if one does not already exist.

## Emergency recovery

```text
sego emergency status
sego emergency unlock
```

Unlock is intentionally a local CLI operation and verifies an administrator password.

## Performance mode

```text
sego mode show
sego mode low-resource
sego mode standard
```

The older explicit syntax is also accepted:

```text
sego mode set low-resource
sego mode set standard
```

## Advanced custom config

Custom deployments can still specify a path explicitly:

```text
sego --config /custom/path/sego.json doctor
```

Normal installations should not need this option.
