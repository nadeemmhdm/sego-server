#!/usr/bin/env bash
set -euo pipefail

PROFILE="auto"
YES=0
PREFIX="${SEGO_PREFIX:-/opt/sego}"
DATA="${SEGO_DATA:-/var/lib/sego}"
CONFIG="${SEGO_CONFIG:-/etc/sego/sego.json}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_VERSION="1.27.2"

usage() {
  cat <<'USAGE'
SEGO Server guided installer (Linux)

Usually you do not need options. Just run the one-command installer from README.

Options:
  --profile auto|low-resource|standard   Default: auto
  --yes                                  Skip install-plan confirmation
  --prefix PATH                          Advanced: binary installation root
  --data PATH                            Advanced: application data directory
  --config PATH                          Advanced: configuration path
  -h, --help                             Show this help
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
case "$PROFILE" in auto|low-resource|standard) ;; *) echo "Invalid profile: $PROFILE" >&2; exit 2 ;; esac
[[ "$(uname -s)" == "Linux" ]] || { echo "This installer is for Linux. Use the Windows PowerShell installer on Windows." >&2; exit 1; }
ARCH="$(uname -m)"
case "$ARCH" in x86_64|amd64) GO_ARCH=amd64; GO_SHA="ecbadb99091a3f46e31f5f934b068b1864eafa7995211b39eaddf76996045fe5" ;; aarch64|arm64) GO_ARCH=arm64; GO_SHA="94f3e30b8e374bc285e7dadc11e0865726b9bc6e85b841ccceaabc0214c6b7c8" ;; *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;; esac

if [[ -r /etc/os-release ]]; then . /etc/os-release; case "${ID:-}" in ubuntu|debian) ;; *) echo "Warning: ${PRETTY_NAME:-this distribution} is outside the initially tested Ubuntu/Debian matrix." ;; esac; fi
if [[ $EUID -eq 0 ]]; then SUDO=(); else command -v sudo >/dev/null 2>&1 || { echo "sudo is required." >&2; exit 1; }; SUDO=(sudo); fi
step(){ printf '\n[%s] %s\n' "$1" "$2"; }

if [[ "$PROFILE" == "auto" ]]; then mem_kb="$(awk '/MemTotal:/ {print $2; exit}' /proc/meminfo 2>/dev/null || true)"; if [[ -n "${mem_kb:-}" && "$mem_kb" -ge 6291456 ]]; then PROFILE="standard"; else PROFILE="low-resource"; fi; fi

missing=()
for p in ca-certificates curl tar nodejs npm; do
  case "$p" in
    ca-certificates) dpkg -s ca-certificates >/dev/null 2>&1 || missing+=(ca-certificates) ;;
    curl) command -v curl >/dev/null 2>&1 || missing+=(curl) ;;
    tar) command -v tar >/dev/null 2>&1 || missing+=(tar) ;;
    nodejs) command -v node >/dev/null 2>&1 || missing+=(nodejs) ;;
    npm) command -v npm >/dev/null 2>&1 || missing+=(npm) ;;
  esac
done
command -v systemctl >/dev/null 2>&1 || { echo "systemd is required by this installer." >&2; exit 1; }
if ((${#missing[@]})); then
  command -v apt-get >/dev/null 2>&1 || { printf 'Missing packages: %s\n' "${missing[*]}" >&2; exit 1; }
  step "1/9" "Installing required system tools"
  if [[ $YES -ne 1 ]]; then read -r -p "Install required packages? [Y/n] " a; [[ "${a:-y}" =~ ^([yY]|[yY][eE][sS])$ ]] || exit 1; fi
  "${SUDO[@]}" apt-get update; "${SUDO[@]}" apt-get install -y "${missing[@]}"
fi

version_ge(){ [[ "$(printf '%s\n%s\n' "$2" "$1" | sort -V | head -n1)" == "$2" ]]; }
ensure_go(){
  local have=""; if command -v go >/dev/null 2>&1; then have="$(go env GOVERSION 2>/dev/null | sed 's/^go//' || true)"; fi
  if [[ -n "$have" ]] && version_ge "$have" "1.24.0"; then echo "Go $have detected."; return; fi
  step "2/9" "Installing verified Go $GO_VERSION toolchain"
  local tmp archive url got; tmp="$(mktemp -d)"; archive="$tmp/go.tar.gz"; url="https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
  curl -fL --retry 3 "$url" -o "$archive"; got="$(sha256sum "$archive" | awk '{print $1}')"
  [[ "$got" == "$GO_SHA" ]] || { rm -rf "$tmp"; echo "Go download checksum verification failed." >&2; exit 1; }
  "${SUDO[@]}" rm -rf /usr/local/go; "${SUDO[@]}" tar -C /usr/local -xzf "$archive"; rm -rf "$tmp"; export PATH="/usr/local/go/bin:$PATH"
}
ensure_go
if ! command -v node >/dev/null 2>&1 || ! command -v npm >/dev/null 2>&1; then echo "Node.js and npm are required." >&2; exit 1; fi
node_major="$(node -p 'process.versions.node.split(".")[0]' 2>/dev/null || echo 0)"; if (( node_major < 18 )); then echo "Node.js 18 or newer is required. Installed: $(node --version)" >&2; exit 1; fi

cat <<PLAN

SEGO is ready to install.
  Mode       : $PROFILE
  Dashboard  : http://127.0.0.1:8088
  Data folder: $DATA
  Startup    : automatic (systemd)
  Updates    : automatic GitHub Release checks

No disk will be formatted or repartitioned. Existing SEGO user data is preserved.
PLAN
if [[ $YES -ne 1 ]]; then read -r -p "Continue? [Y/n] " a; [[ "${a:-y}" =~ ^([yY]|[yY][eE][sS])$ ]] || exit 0; fi

TMP="$(mktemp -d)"; trap 'rm -rf "$TMP"' EXIT; cd "$REPO_ROOT"
step "3/9" "Resolving verified Go dependencies"; go mod tidy
step "4/9" "Building the dashboard"; (cd web; npm install --no-audit --no-fund; npm run build)
step "5/9" "Building SEGO"
go build -trimpath -o "$TMP/sego-server" ./cmd/server
go build -trimpath -o "$TMP/sego" ./cmd/sego
go build -trimpath -o "$TMP/sego-updater" ./cmd/updater

step "6/9" "Creating secure folders"
if ! id -u sego >/dev/null 2>&1; then "${SUDO[@]}" useradd --system --home "$DATA" --shell /usr/sbin/nologin sego; fi
"${SUDO[@]}" install -d -m 0755 "$PREFIX/bin" "$(dirname "$CONFIG")"
"${SUDO[@]}" install -d -o sego -g sego -m 0750 "$DATA" "$DATA/web"

step "7/9" "Installing SEGO, boot startup and automatic updates"
"${SUDO[@]}" install -m 0755 "$TMP/sego-server" "$PREFIX/bin/sego-server"
"${SUDO[@]}" install -m 0755 "$TMP/sego" "$PREFIX/bin/sego"
"${SUDO[@]}" install -m 0755 "$TMP/sego-updater" "$PREFIX/bin/sego-updater"
"${SUDO[@]}" ln -sf "$PREFIX/bin/sego" /usr/local/bin/sego
"${SUDO[@]}" rm -rf "$DATA/web"; "${SUDO[@]}" install -d -o sego -g sego -m 0750 "$DATA/web"; "${SUDO[@]}" cp -a web/dist/. "$DATA/web/"; "${SUDO[@]}" chown -R sego:sego "$DATA/web"
"${SUDO[@]}" env SEGO_DATA_DIR="$DATA" "$PREFIX/bin/sego" --config "$CONFIG" mode "$PROFILE"
"${SUDO[@]}" chown -R sego:sego "$DATA"; "${SUDO[@]}" chown root:sego "$CONFIG"; "${SUDO[@]}" chmod 0640 "$CONFIG"

sed -e "s#/opt/sego/bin/sego-server#$PREFIX/bin/sego-server#g" -e "s#/etc/sego/sego.json#$CONFIG#g" -e "s#/var/lib/sego#$DATA#g" deploy/systemd/sego.service > "$TMP/sego.service"
sed -e "s#/opt/sego#$PREFIX#g" -e "s#/etc/sego/sego.json#$CONFIG#g" -e "s#/var/lib/sego#$DATA#g" deploy/systemd/sego-update.service > "$TMP/sego-update.service"
"${SUDO[@]}" install -m 0644 "$TMP/sego.service" /etc/systemd/system/sego.service
"${SUDO[@]}" install -m 0644 "$TMP/sego-update.service" /etc/systemd/system/sego-update.service
"${SUDO[@]}" install -m 0644 deploy/systemd/sego-update.timer /etc/systemd/system/sego-update.timer
"${SUDO[@]}" systemctl daemon-reload
"${SUDO[@]}" systemctl enable --now sego.service
"${SUDO[@]}" systemctl enable --now sego-update.timer

step "8/9" "Creating your administrator account"
"${SUDO[@]}" "$PREFIX/bin/sego" --config "$CONFIG" doctor
"${SUDO[@]}" systemctl --no-pager --full is-active sego.service >/dev/null || { echo "SEGO could not start. Run: sudo journalctl -u sego.service -n 100 --no-pager" >&2; exit 1; }
if [[ -t 0 ]]; then "${SUDO[@]}" "$PREFIX/bin/sego" --config "$CONFIG" setup; else echo "No interactive terminal is available. After install run: sudo sego setup"; fi

step "9/9" "Installation complete"
printf '\nOpen SEGO:  sego open\nCheck SEGO: sego status\nRemote:     sego tunnel setup\nUpdates:    sego update check\nDashboard:  http://127.0.0.1:8088\n'
if [[ -n "${DISPLAY:-}${WAYLAND_DISPLAY:-}" ]] && command -v xdg-open >/dev/null 2>&1; then xdg-open http://127.0.0.1:8088 >/dev/null 2>&1 || true; fi
