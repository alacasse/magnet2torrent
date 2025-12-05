#!/usr/bin/env bash

# Register magnet2torrent as the handler for magnet: links on Linux.
# Writes a desktop entry under ~/.local/share/applications and sets x-scheme-handler/magnet via xdg-mime.

set -euo pipefail

APP_NAME="${APP_NAME:-magnet2torrent}"
DESKTOP_FILE="${DESKTOP_FILE:-$HOME/.local/share/applications/${APP_NAME}.desktop}"
DRY_RUN="${DRY_RUN:-0}"

usage() {
  cat <<EOF
Usage: APP_NAME=magnet2torrent DESKTOP_FILE=\$HOME/.local/share/applications/magnet2torrent.desktop $(basename "$0") [--dry-run]

Env vars:
  APP_NAME      Binary name to register (default: magnet2torrent)
  DESKTOP_FILE  Full path to desktop file (default: ~/.local/share/applications/magnet2torrent.desktop)
  DRY_RUN       Set to 1 to print actions without executing
EOF
}

run() {
  if [[ "$DRY_RUN" == "1" ]]; then
    echo "[dry-run] $*"
  else
    eval "$@"
  fi
}

resolve_binary() {
  local bin
  if ! bin="$(command -v "$APP_NAME")"; then
    echo "Could not find $APP_NAME in PATH; install it first." >&2
    exit 1
  fi
  echo "$bin"
}

write_desktop_file() {
  local bin_path="$1"
  local dir
  dir="$(dirname "$DESKTOP_FILE")"

  read -r -d '' desktop_content <<EOF || true
[Desktop Entry]
Type=Application
Name=${APP_NAME}
Comment=Send magnet links to magnet2torrent
Exec=${bin_path} %u
NoDisplay=true
Terminal=false
MimeType=x-scheme-handler/magnet;
Categories=Network;FileTransfer;
EOF

  if [[ "$DRY_RUN" == "1" ]]; then
    echo "[dry-run] write $DESKTOP_FILE with:"
    echo "$desktop_content"
    return
  fi

  run "mkdir -p \"$dir\""
  printf '%s\n' "$desktop_content" > "$DESKTOP_FILE"
  echo "Wrote $DESKTOP_FILE"
}

register_mime() {
  if ! command -v xdg-mime >/dev/null 2>&1; then
    echo "xdg-mime not found; cannot set magnet handler automatically." >&2
    exit 1
  fi
  run "xdg-mime default \"$(basename "$DESKTOP_FILE")\" x-scheme-handler/magnet"

  if command -v update-desktop-database >/dev/null 2>&1; then
    run "update-desktop-database \"$(dirname "$DESKTOP_FILE")\""
  fi
}

main() {
  if [[ "${1:-}" == "--help" ]]; then
    usage
    exit 0
  fi
  if [[ "${1:-}" == "--dry-run" ]]; then
    DRY_RUN=1
  fi

  local bin_path
  bin_path="$(resolve_binary)"
  write_desktop_file "$bin_path"
  register_mime

  echo "Registered magnet: handler for $APP_NAME"
}

main "$@"
