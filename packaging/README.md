# Packaging

This folder tracks packaging assets for distribution channels.

## Arch (AUR)

Use the -git package for now. Build on Arch with:

```bash
cd packaging/aur/bubblefetch-git
makepkg -si
```

Maintainer: Howie Duhzit <Contact@HowieDuhzit.Best>

Note: `pkgver` is computed from git tags and commit count. Run:

```bash
makepkg --printsrcinfo > .SRCINFO
```

after any PKGBUILD edits.
