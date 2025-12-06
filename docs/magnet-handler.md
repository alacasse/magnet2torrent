# Magnet Link Handler

The app can register itself as the handler for `magnet:` links.

## Linux

Desktops (GNOME, KDE/Plasma, XFCE, Cinnamon, etc.) rely on a desktop entry plus `xdg-mime`.

**Register**

```bash
scripts/register-magnet-linux.sh
```

This writes `~/.local/share/applications/magnet2torrent.desktop`, then runs `xdg-mime default magnet2torrent.desktop x-scheme-handler/magnet` and `update-desktop-database ~/.local/share/applications` when available. Use `DRY_RUN=1` to preview changes.

The installer (`scripts/install.sh`) and the npm package postinstall call the register script automatically on Linux. Set `REGISTER_MAGNET=0` to skip.

**Unregister**

```bash
xdg-mime default '' x-scheme-handler/magnet && \
rm -f ~/.local/share/applications/magnet2torrent.desktop && \
update-desktop-database ~/.local/share/applications 2>/dev/null || true
```

You may need to reselect a different app in your browser after removing the handler.

## Windows

Windows uses per-user registry keys under `HKCU:\Software\Classes\magnet`.

**Register**

```powershell
pwsh scripts/register-magnet-windows.ps1
```

This sets the magnet protocol command to the `magnet2torrent` binary for the current user. Use `-DryRun` to preview changes. The Windows installer (`scripts/install.ps1`) and the npm package postinstall invoke this automatically unless `REGISTER_MAGNET=0` or `-RegisterMagnet:$false` is provided.

If Windows prompts for a handler after registration, open Settings → Apps → Default apps → Choose defaults by link type, search for `magnet`, and pick `magnet2torrent`.

**Unregister**

```powershell
Remove-Item -Path HKCU:\Software\Classes\magnet -Recurse -Force
```

You may need to reassign magnet links to another app in Default apps afterward.
