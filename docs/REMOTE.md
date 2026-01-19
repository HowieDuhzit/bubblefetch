# Remote Collection (SSH)

Bubblefetch can collect system info over SSH using `--remote` (or config `remote:`).
It uses your local SSH config and authentication (keys or agent).

## What Runs on the Remote Host

Default mode uses a shell wrapper to run simple commands and parse their output:

- `/etc/os-release` or `/usr/lib/os-release`
- `uname -r`
- `hostname` or `/etc/hostname`
- `/proc/uptime`
- `/proc/cpuinfo`
- `/proc/meminfo`
- `df -B1 /`
- `lspci` (GPU)
- `ip -o addr show` (network)
- `/sys/class/power_supply/BAT*/capacity` (battery)
- environment vars: `$SHELL`, `$TERM`, `$XDG_CURRENT_DESKTOP`, `$XDG_SESSION_TYPE`

The wrapper sets `PATH` and `LC_ALL=C` for consistent output.

## Safe Mode (Read-Only)

Use `--remote-safe` or set `ssh.safe_mode: true` to avoid shell pipelines and
prefer read-only file access. Safe mode currently collects:

- OS, kernel, hostname, uptime
- CPU (from `/proc/cpuinfo`)
- Memory (from `/proc/meminfo`)

Disk, GPU, network, shell, terminal, and DE/WM are omitted in safe mode.

## Permissions & Security

Remote commands run as the SSH user you connect with. Use least-privilege accounts
and restrict SSH access as appropriate for your environment.
