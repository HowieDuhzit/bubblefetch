# Quick Start Guide

Get bubblefetch up and running in 60 seconds!

Need more detail? Start at `docs/README.md`.

## Installation

```bash
git clone https://github.com/howieduhzit/bubblefetch.git
cd bubblefetch
./install.sh
```

That's it! The script handles everything:
- Builds the binary
- Installs to `/usr/local/bin`
- Sets up config directory
- Copies themes

## First Run

```bash
bubblefetch
```

You should see your system info with a beautiful ASCII art logo!

## Try Different Themes

```bash
bubblefetch --theme dracula
bubblefetch --theme gruvbox
bubblefetch --theme tokyo-night
bubblefetch --theme nord
bubblefetch --theme minimal
```

## Customize

Edit your config:

```bash
nano ~/.config/bubblefetch/config.yaml
```

Example config:

```yaml
theme: gruvbox

modules:
  - os
  - kernel
  - uptime
  - cpu
  - gpu
  - memory
  - disk
  - localip
```

Then run:

```bash
bubblefetch
```

## Export Your System Info

```bash
# JSON format
bubblefetch --export json > my-system.json

# YAML format
bubblefetch --export yaml > my-system.yaml

# Plain text
bubblefetch --export text > my-system.txt
```

## Check a Remote Server

```bash
bubblefetch --remote user@hostname
```

Uses your existing SSH keys - no setup needed!

## Benchmark Performance

```bash
bubblefetch --benchmark
```

Runs 10 iterations and shows average collection time.

## Next Steps

- Read [README.md](README.md) for full documentation
- Check [EXAMPLES.md](EXAMPLES.md) for advanced usage
- Browse themes in `~/.config/bubblefetch/themes/`
- Create custom themes (see README)

## Common Issues

**"command not found"**
- Make sure `/usr/local/bin` is in your PATH
- Try `export PATH=$PATH:/usr/local/bin`

**No GPU detected**
- Some systems may not have `lspci` installed
- Install with: `sudo apt install pciutils` (Debian/Ubuntu)

**SSH connection fails**
- Check your SSH keys: `ls ~/.ssh/`
- Test SSH manually: `ssh user@host`
- Verify SSH config if using host aliases

## Uninstall

```bash
cd bubblefetch
./uninstall.sh
```

## Get Help

- GitHub Issues: https://github.com/howieduhzit/bubblefetch/issues
- Run: `bubblefetch --help`
