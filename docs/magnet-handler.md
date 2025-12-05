# Magnet Link Handler (Linux)

This app can register itself as the handler for `magnet:` links on Linux desktops (GNOME, KDE/Plasma, XFCE, Cinnamon, etc.) by installing a desktop entry and telling `xdg-mime` to route the protocol to `magnet2torrent`.

## Register

After installing the CLI, run:

```bash
scripts/register-magnet-linux.sh
```

By default this writes `~/.local/share/applications/magnet2torrent.desktop`, then runs `xdg-mime default magnet2torrent.desktop x-scheme-handler/magnet` and `update-desktop-database ~/.local/share/applications` when available. Use `DRY_RUN=1` to preview changes.

The installer (`scripts/install.sh`) calls the register script automatically on Linux. Set `REGISTER_MAGNET=0` to skip.

## Unregister

```bash
xdg-mime default '' x-scheme-handler/magnet && \
rm -f ~/.local/share/applications/magnet2torrent.desktop && \
update-desktop-database ~/.local/share/applications 2>/dev/null || true
```

You may need to reselect a different app in your browser after removing the handler.
