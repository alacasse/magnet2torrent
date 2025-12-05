#!/usr/bin/env bash

# Extensible installer for magnet2torrent via npm tarball from GitHub Releases.
# Supports linux and macOS today; extend per-OS functions to add config prompts or alternate installers.

set -euo pipefail

REPO_OWNER="${REPO_OWNER:-alacasse}"
REPO_NAME="${REPO_NAME:-magnet2torrent}"
VERSION="${VERSION:-}"
PREFIX="${PREFIX:-}"
DRY_RUN="${DRY_RUN:-0}"
REGISTER_MAGNET="${REGISTER_MAGNET:-1}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

usage() {
  cat <<EOF
Usage: REPO_OWNER=you REPO_NAME=repo VERSION=0.1.0 $(basename "$0") [--prefix /custom/npm/prefix] [--dry-run]

Env vars:
  REPO_OWNER   GitHub owner/org (default: you)
  REPO_NAME    GitHub repo name (default: repo)
  VERSION      Release version without v prefix (required)
  PREFIX       npm prefix to install under (optional)
  DRY_RUN      Set to 1 to print actions without executing
  REGISTER_MAGNET Set to 0 to skip magnet handler registration on Linux
EOF
}

run() {
  if [[ "$DRY_RUN" == "1" ]]; then
    echo "[dry-run] $*"
  else
    eval "$@"
  fi
}

ensure_version() {
  if [[ -z "$VERSION" ]]; then
    echo "VERSION is required (e.g., VERSION=0.1.0)." >&2
    usage
    exit 1
  fi
}

tarball_url() {
  echo "https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/v$VERSION/magnet2torrent-$VERSION.tgz"
}

pre_install_config() {
  # Hook for future interactive config collection.
  :
}

register_magnet_linux() {
  if [[ "$REGISTER_MAGNET" != "1" ]]; then
    echo "Skipping magnet handler registration (REGISTER_MAGNET=$REGISTER_MAGNET)"
    return
  fi
  if [[ ! -x "$SCRIPT_DIR/register-magnet-linux.sh" ]]; then
    echo "Registration script not found at $SCRIPT_DIR/register-magnet-linux.sh" >&2
    return
  fi
  echo "Registering magnet handler..."
  run "DRY_RUN=\"$DRY_RUN\" \"$SCRIPT_DIR/register-magnet-linux.sh\""
}

install_linux() {
  local url
  url="$(tarball_url)"
  echo "Installing magnet2torrent $VERSION for linux (npm global)..."
  if [[ -n "$PREFIX" ]]; then
    run "npm install -g --prefix \"$PREFIX\" \"$url\""
  else
    run "npm install -g \"$url\""
  fi
  register_magnet_linux
}

install_macos() {
  local url
  url="$(tarball_url)"
  echo "Installing magnet2torrent $VERSION for macOS (npm global)..."
  if [[ -n "$PREFIX" ]]; then
    run "npm install -g --prefix \"$PREFIX\" \"$url\""
  else
    run "npm install -g \"$url\""
  fi
}

post_install_message() {
  echo "magnet2torrent installed. Verify with: magnet2torrent --help"
  if [[ -n "$PREFIX" ]]; then
    echo "If not on PATH, add: export PATH=\"$PREFIX/bin:\$PATH\""
  fi
  if [[ "$(uname -s)" == "Linux" && "$REGISTER_MAGNET" != "1" ]]; then
    echo "Magnet handler registration skipped (REGISTER_MAGNET=$REGISTER_MAGNET). Run scripts/register-magnet-linux.sh later if needed."
  fi
}

main() {
  ensure_version
  pre_install_config

  case "$(uname -s)" in
    Linux) install_linux ;;
    Darwin) install_macos ;;
    *)
      echo "Unsupported OS: $(uname -s)" >&2
      exit 1
      ;;
  esac

  post_install_message
}

main "$@"
