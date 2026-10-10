#!/usr/bin/env bash
set -euo pipefail
PREFIX="${SEGO_PREFIX:-/opt/sego}"
CONFIG="${SEGO_CONFIG:-/etc/sego/sego.json}"
DATA="${SEGO_DATA:-/var/lib/sego}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
[[ -x /usr/local/go/bin/go ]] && export PATH="/usr/local/go/bin:$PATH"
if [[ $EUID -eq 0 ]]; then SUDO=(); else SUDO=(sudo); fi
cd "$ROOT"
printf '\nSEGO Configuring automatic release updates\n'
tmp="$(mktemp -d)"; trap 'rm -rf "$tmp"' EXIT
go build -trimpath -o "$tmp/sego-updater" ./cmd/updater
"${SUDO[@]}" install -m 0755 "$tmp/sego-updater" "$PREFIX/bin/sego-updater"
sed -e "s#/opt/sego#$PREFIX#g" -e "s#/etc/sego/sego.json#$CONFIG#g" -e "s#/var/lib/sego#$DATA#g" deploy/systemd/sego-update.service > "$tmp/sego-update.service"
"${SUDO[@]}" install -m 0644 "$tmp/sego-update.service" /etc/systemd/system/sego-update.service
"${SUDO[@]}" install -m 0644 deploy/systemd/sego-update.timer /etc/systemd/system/sego-update.timer
"${SUDO[@]}" systemctl daemon-reload
"${SUDO[@]}" systemctl enable --now sego-update.timer
printf '[OK] Automatic GitHub Release checks enabled.\n'
