#!/usr/bin/env bash
# SEGO Server one-command beginner installer.
# curl -fsSL https://raw.githubusercontent.com/nadeemmhdm/sego-server/main/install.sh | bash
set -euo pipefail
REPO="nadeemmhdm/sego-server";PROFILE="${SEGO_PROFILE:-auto}";TMP="$(mktemp -d)";trap 'rm -rf "$TMP"' EXIT
say(){ printf '\n\033[1;36mSEGO\033[0m %s\n' "$*"; };fail(){ printf '\nSEGO install error: %s\n' "$*" >&2;exit 1; }
command -v curl>/dev/null||fail "curl is required: sudo apt install curl";command -v tar>/dev/null||fail "tar is required: sudo apt install tar"
say "Server installer";echo "SEGO will download itself from GitHub, install required tools, configure automatic startup and guide first setup."
if [[ -r /dev/tty ]];then printf 'Continue? [Y/n] '>/dev/tty;read -r a</dev/tty||true;[[ "${a:-y}" =~ ^([nN]|[nN][oO])$ ]]&&exit 0;fi
say "Downloading SEGO from GitHub";meta="$(curl -fsSL -H 'Accept: application/vnd.github+json' "https://api.github.com/repos/$REPO/commits/main")"||fail "Could not contact GitHub";commit="$(printf '%s' "$meta"|sed -n 's/.*"sha"[[:space:]]*:[[:space:]]*"\([0-9a-f]\{40\}\)".*/\1/p'|head -n1)";[[ "$commit" =~ ^[0-9a-f]{40}$ ]]||fail "Could not resolve main commit";curl -fL --retry 3 "https://github.com/$REPO/archive/$commit.tar.gz" -o "$TMP/sego.tar.gz";tar -xzf "$TMP/sego.tar.gz" -C "$TMP";ROOT="$TMP/sego-server-$commit"
if [[ -r /dev/tty ]];then bash "$ROOT/scripts/install.sh" --profile "$PROFILE" --yes </dev/tty;else bash "$ROOT/scripts/install.sh" --profile "$PROFILE" --yes;fi
bash "$ROOT/scripts/install-update-service.sh"
echo;echo "SEGO starts automatically whenever this computer boots.";echo "Automatic GitHub Release update checks are enabled.";echo "Optional remote access: run  sego tunnel setup";echo "Open SEGO anytime with: sego open"
if [[ -r /dev/tty ]];then printf 'Set up Cloudflare Tunnel remote access now? [y/N] '>/dev/tty;read -r t</dev/tty||true;if [[ "$t" =~ ^([yY]|[yY][eE][sS])$ ]];then /usr/local/bin/sego tunnel setup </dev/tty;fi;fi
