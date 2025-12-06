# magnet2torrent

CLI helper that forwards `magnet:` links to qBittorrent via its WebUI API. It prompts for (or reads) your qBittorrent host/username/password, then accepts a magnet link and hands it off to qBittorrent.

## Install

### Prerequisites

- Node.js + npm (for the provided installers)
- qBittorrent with WebUI enabled and reachable from this machine

### Linux

#### Install from GitHub release (recommended)

```bash
# replace VERSION with the release you want, e.g. 0.1.0
VERSION=0.1.0 scripts/install.sh
```

This installs the npm tarball globally and registers `magnet2torrent` as the default handler for `magnet:` links. Skip handler registration with `REGISTER_MAGNET=0` and register later via `scripts/register-magnet-linux.sh`.

#### Install from a locally downloaded tarball

If you downloaded `magnet2torrent-VERSION.tgz` already:

```bash
REGISTER_MAGNET=0 npm install -g /path/to/magnet2torrent-0.1.0.tgz
# registration now happens automatically unless REGISTER_MAGNET=0
```

#### Build from source

```bash
make build
sudo cp bin/magnet2torrent /usr/local/bin/
# optional: register magnet handler
scripts/register-magnet-linux.sh
```

### Windows

#### Install from GitHub release (recommended)

```powershell
powershell -ExecutionPolicy Bypass -File scripts/install.ps1 -Version 0.1.0
```

This installs the npm tarball globally and registers the `magnet:` protocol to `magnet2torrent` for the current user. Skip registration with `-RegisterMagnet:$false` or `REGISTER_MAGNET=0`, then run `scripts/register-magnet-windows.ps1` later.

#### Install from a locally downloaded tarball

If you downloaded `magnet2torrent-VERSION.tgz` already:

```powershell
REGISTER_MAGNET=0 npm install -g C:\path\to\magnet2torrent-0.1.0.tgz
# registration now happens automatically unless REGISTER_MAGNET=0
```

#### Build from source

```powershell
go build -o bin\magnet2torrent .\cmd\magnet2torrent
Copy-Item bin\magnet2torrent.exe $env:USERPROFILE\AppData\Local\Microsoft\WindowsApps\
powershell -ExecutionPolicy Bypass -File scripts/register-magnet-windows.ps1
```

## Configuration

First run will prompt for:

- `qbHost` (e.g., `http://localhost:8080`)
- `qbUsername`
- `qbPassword`

Config is stored at `~/.config/magnet2torrent/config.json` (Linux) or `%APPDATA%\magnet2torrent\config.json` (Windows). Edit or pre-create it to skip prompts.

### Logging

Logs go to stdout and to the log file defined in config (`logFile`), defaulting to `~/.cache/magnet2torrent/magnet2torrent.log` on Linux and `%LOCALAPPDATA%\magnet2torrent\magnet2torrent.log` on Windows. Use this file to inspect runs triggered via browser magnet links.

## Usage

```bash
magnet2torrent "magnet:?xt=urn:btih:..."
```

Flags:

- `-config <path>`: path to a config file
- `-v` / `-version`: print version and exit

## Magnet handler registration

- Linux: `scripts/register-magnet-linux.sh` writes a desktop entry to `~/.local/share/applications` and calls `xdg-mime default magnet2torrent.desktop x-scheme-handler/magnet`.
- Windows: `scripts/register-magnet-windows.ps1` sets `HKCU:\Software\Classes\magnet` to point to `magnet2torrent`.

If your browser prompts after registration, choose magnet2torrent and allow it to remember the choice. Unregister steps are documented in `docs/magnet-handler.md`.
