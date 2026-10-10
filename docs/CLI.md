# SEGO CLI

SEGO keeps common commands short. Normal users do not need to know configuration paths.

## Beginner commands

```text
sego             Start SEGO if needed and open the dashboard
sego status      Show server/service status
sego update      Run the installed release updater
sego remote      Start the guided Cloudflare Tunnel setup
sego unlock      Local Emergency SOS recovery
sego help        Show the complete command list
```

Service controls:

```text
sego start
sego stop
sego restart
```

Performance shortcuts:

```text
sego low         Low Resource Mode
sego normal      Standard Mode
```

`sego check` is an alias for `sego status`.

## Administrator / compatibility commands

The explicit commands remain available for scripts and administrators:

```text
sego open
sego setup
sego doctor
sego storage list
sego tunnel setup
sego tunnel status
sego tunnel remove
sego update check
sego update auto
sego emergency status
sego emergency unlock
sego mode show
sego mode low-resource
sego mode standard
```

On Windows, service-control commands can show a UAC prompt. On Linux, sudo can request the user's password when the operating system requires elevation.

## Advanced custom configuration

Custom deployments can still use:

```text
sego --config /custom/path/sego.json doctor
```

Normal installations should not need this option.
