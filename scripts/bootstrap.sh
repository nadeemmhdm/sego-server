#!/usr/bin/env bash
# SEGO Server beginner installer bootstrap.
# Usage: curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
set -euo pipefail

REPO="nadeemmhdm/sego-server"
PROFILE="${SEGO_PROFILE:-auto}"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

say(){ printf '\n\033[1;36mSEGO\033[0m %s\n' "$*"; }
fail(){ printf '\nSEGO install error: %s\n' "$*" >&2; exit 1; }

command -v curl >/dev/null 2>&1 || fail "curl is required. On Ubuntu/Debian run: sudo apt install curl"
command -v tar >/dev/null 2>&1 || fail "tar is required. On Ubuntu/Debian run: sudo apt install tar"

say "Server installer"
echo "This will automatically download SEGO from GitHub, install required tools,"
echo "create the background service, and guide you through the first admin account."

if [[ -r /dev/tty ]]; then
  printf 'Continue? [Y/n] ' >/dev/tty
  read -r answer </dev/tty || true
  case "${answer:-y}" in n|N|no|NO|No) echo "Installation cancelled."; exit 0;; esac
fi

say "Downloading the current GitHub version"
meta="$(curl -fsSL -H 'Accept: application/vnd.github+json' "https://api.github.com/repos/$REPO/commits/main")" || fail "Could not contact GitHub."
commit="$(printf '%s' "$meta" | sed -n 's/.*"sha"[[:space:]]*:[[:space:]]*"\([0-9a-f]\{40\}\)".*/\1/p' | head -n1)"
[[ "$commit" =~ ^[0-9a-f]{40}$ ]] || fail "Could not resolve the GitHub main commit."
curl -fL --retry 3 "https://github.com/$REPO/archive/$commit.tar.gz" -o "$TMP/sego.tar.gz" || fail "SEGO source download failed."
tar -xzf "$TMP/sego.tar.gz" -C "$TMP"
installer="$TMP/sego-server-$commit/scripts/install.sh"
[[ -x "$installer" || -f "$installer" ]] || fail "Downloaded installer is missing."

say "Starting guided setup"
if [[ -r /dev/tty ]]; then
  bash "$installer" --profile "$PROFILE" --yes </dev/tty
else
  bash "$installer" --profile "$PROFILE" --yes
fi
