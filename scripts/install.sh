#!/usr/bin/env bash
set -euo pipefail

PROFILE="low-resource"
YES=0
PREFIX="${SEGO_PREFIX:-/opt/sego}"
DATA="${SEGO_DATA:-/var/lib/sego}"
CONFIG="${SEGO_CONFIG:-/etc/sego/sego.json}"

usage() {
  cat <<'USAGE'
SEGO Server installer (Linux)

Usage:
  ./scripts/install.sh [--profile low-resource|standard] [--yes]

Options:
  --profile NAME   Runtime profile. Default: low-resource
  --prefix PATH    Binary installation root. Default: /opt/sego
  --data PATH      Application data directory. Default: /var/lib/sego
  --config PATH    Configuration path. Default: /etc/sego/sego.json
  --yes            Non-interactive approval of the displayed install plan
  -h, --help       Show this help
USAGE
}

while (($#)); do
  case "$1" in
    --profile) PROFILE="${2:?missing profile}"; shift 2 ;;
    --prefix) PREFIX="${2:?missing prefix}"; shift 2 ;;
    --data) DATA="${2:?missing data path}"; shift 2 ;;
    --config) CONFIG="${2:?missing config path}"; shift 2 ;;
    --yes|-y) YES=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

case "$PROFILE" in
  low-resource|standard) ;;
  *) echo "Invalid profile: $PROFILE" >&2; exit 2 ;;
esac

[[ "$(uname -s)" == "Linux" ]] || { echo "This installer is for Linux. Use scripts/install.ps1 on Windows." >&2; exit 1; }
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64|amd64|aarch64|arm64) ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

if [[ -r /etc/os-release ]]; then
  # shellcheck disable=SC1091
  . /etc/os-release
  case "${ID:-}" in
    ubuntu|debian) ;;
    *) echo "Warning: ${PRETTY_NAME:-this Linux distribution} is not in the initially tested Ubuntu/Debian support matrix." ;;
  esac
fi

command -v go >/dev/null 2>&1 || { echo "Go 1.23+ is required when installing from a source checkout." >&2; exit 1; }
command -v systemctl >/dev/null 2>&1 || { echo "systemd is required by this source installer." >&2; exit 1; }

if [[ $EUID -eq 0 ]]; then
  SUDO=()
else
  command -v sudo >/dev/null 2>&1 || { echo "sudo is required for system installation." >&2; exit 1; }
  SUDO=(sudo)
fi

cat <<PLAN
SEGO Server install plan
  profile : $PROFILE
  binaries: $PREFIX/bin
  data    : $DATA
  config  : $CONFIG
  service : systemd (sego.service)

The installer never formats or repartitions disks. Existing user data is preserved.
PLAN

if [[ $YES -ne 1 ]]; then
  read -r -p "Continue? [y/N] " answer
  [[ "${answer,,}" == "y" || "${answer,,}" == "yes" ]] || exit 0
fi

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "[1/7] Building SEGO binaries"
go build -trimpath -o "$TMP/sego-server" ./cmd/server
go build -trimpath -o "$TMP/sego" ./cmd/sego

echo "[2/7] Creating service account and directories"
if ! id -u sego >/dev/null 2>&1; then
  "${SUDO[@]}" useradd --system --home "$DATA" --shell /usr/sbin/nologin sego
fi
"${SUDO[@]}" install -d -m 0755 "$PREFIX/bin" "$(dirname "$CONFIG")"
"${SUDO[@]}" install -d -o sego -g sego -m 0750 "$DATA"

echo "[3/7] Installing binaries"
"${SUDO[@]}" install -m 0755 "$TMP/sego-server" "$PREFIX/bin/sego-server"
"${SUDO[@]}" install -m 0755 "$TMP/sego" "$PREFIX/bin/sego"

echo "[4/7] Creating configuration and applying profile"
"${SUDO[@]}" "$PREFIX/bin/sego" --config "$CONFIG" mode set "$PROFILE"
"${SUDO[@]}" chown -R sego:sego "$DATA"
"${SUDO[@]}" chown root:sego "$CONFIG"
"${SUDO[@]}" chmod 0640 "$CONFIG"

echo "[5/7] Installing systemd service"
SERVICE_TMP="$TMP/sego.service"
sed \
  -e "s#/opt/sego/bin/sego-server#$PREFIX/bin/sego-server#g" \
  -e "s#/etc/sego/sego.json#$CONFIG#g" \
  -e "s#/var/lib/sego#$DATA#g" \
  deploy/systemd/sego.service > "$SERVICE_TMP"
"${SUDO[@]}" install -m 0644 "$SERVICE_TMP" /etc/systemd/system/sego.service
"${SUDO[@]}" systemctl daemon-reload
"${SUDO[@]}" systemctl enable --now sego.service

echo "[6/7] Running health diagnostics"
"${SUDO[@]}" "$PREFIX/bin/sego" --config "$CONFIG" doctor
"${SUDO[@]}" systemctl --no-pager --full status sego.service >/dev/null || {
  echo "SEGO service did not start successfully. Run: sudo journalctl -u sego.service -n 100 --no-pager" >&2
  exit 1
}

echo "[7/7] Installation complete"
printf '\nCreate the first administrator:\n  sudo %q --config %q admin create\n' "$PREFIX/bin/sego" "$CONFIG"
printf 'Dashboard/API (local): http://127.0.0.1:8088\n'
printf 'Change resource mode later:\n  sudo %q --config %q mode set low-resource\n  sudo %q --config %q mode set standard\n' "$PREFIX/bin/sego" "$CONFIG" "$PREFIX/bin/sego" "$CONFIG"
