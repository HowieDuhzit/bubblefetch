# Packaging

This folder tracks packaging assets for distribution channels.

## Arch (AUR)

Packages available:

- `bubblefetch` - release tarball build
- `bubblefetch-bin` - prebuilt release binary
- `bubblefetch-git` - latest git master

Build on Arch with:

```bash
cd packaging/aur/bubblefetch
makepkg -si

cd packaging/aur/bubblefetch-bin
makepkg -si

cd packaging/aur/bubblefetch-git
makepkg -si
```
