#!/usr/bin/env bash
set -euo pipefail
PURGE_DATA=0
PREFIX="${SEGO_PREFIX:-/opt/sego}"
DATA="${SEGO_DATA:-/var/lib/sego}"
CONFIG_DIR="$(dirname "${SEGO_CONFIG:-/etc/sego/sego.json}")"
for arg in "$@"; do
  case "$arg" in
    --purge-data) PURGE_DATA=1 ;;
    -h|--help) echo "Usage: $0 [--purge-data]"; exit 0 ;;
    *) echo "Unknown option: $arg" >&2; exit 2 ;;
  esac
done
SUDO=(); [[ $EUID -eq 0 ]] || SUDO=(sudo)
"${SUDO[@]}" systemctl disable --now sego.service 2>/dev/null || true
"${SUDO[@]}" rm -f /etc/systemd/system/sego.service
"${SUDO[@]}" systemctl daemon-reload
"${SUDO[@]}" rm -rf "$PREFIX"
if [[ $PURGE_DATA -eq 1 ]]; then
  echo "Purging SEGO data and configuration: $DATA $CONFIG_DIR"
  "${SUDO[@]}" rm -rf "$DATA" "$CONFIG_DIR"
else
  echo "SEGO binaries removed. User data and configuration were preserved."
fi
