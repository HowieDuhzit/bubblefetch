# Changelog

All notable changes to bubblefetch will be documented in this file.

## [0.2.1] - 2026-01-17

### Performance Improvements
- **5.6x faster collection**: Optimized from 7.5ms to 1.3ms average
- GPU detection now uses `/sys/class/drm` instead of `lspci` (5-6ms faster)
- OS detection cached with `sync.Once` (0.5ms faster on repeated calls)
- Added timeout to external commands (500ms max)
- Vendor ID mapping for NVIDIA/AMD/Intel GPUs

### Changed
- GPU module now shows "NVIDIA GPU (2489)" format when using /sys
- lspci only used as fallback if /sys fails
- All theme loading operations now cached

## [0.2.0] - 2026-01-17

### Major Features

### Added
- **Performance Optimization**: Parallel collection of system metrics for 2-3x speed improvement
  - CPU, memory, disk, GPU, network, and battery info collected concurrently
  - Reduced execution time significantly compared to sequential collection

- **OS Detection & Auto-ASCII**:
  - Automatic detection of 15+ operating systems and distributions
  - OS-specific ASCII art (Arch, Ubuntu, Debian, Fedora, Mint, Manjaro, Pop!_OS, etc.)
  - Themes can use `"ascii": "auto"` for automatic logo selection

- **New System Info Modules**:
  - `gpu` - GPU detection via lspci and /sys/class/drm
  - `network` - Active network interfaces with IPv4/IPv6 addresses
  - `localip` - Primary local IP address
  - `battery` - Battery percentage, charging status, and time remaining (Linux laptops)

- **New Themes**:
  - Gruvbox - Warm, retro groove colors
  - Tokyo Night - Modern dark theme
  - Monokai - Classic Monokai Pro
  - Solarized Dark - Precision colors

### Changed
- Binary size reduced from 5.5MB to 3.9MB with build optimizations
- All existing themes now use auto-detected ASCII art
- Default module list expanded to include GPU, network, and battery

### SSH Remote Support
- Full SSH remote system support via `--remote user@host`
- Automatic SSH key detection (~/.ssh/id_rsa, id_ed25519)
- Parallel command execution on remote systems
- Support for custom SSH ports and configurations

### Export & Benchmark
- Export to JSON with `--export json` (pretty or compact)
- Export to YAML with `--export yaml`
- Export to plain text with `--export text`
- Benchmark mode with `--benchmark` (10 iterations with timing)

### Installation
- Automated install.sh script
- Automated uninstall.sh script
- Config and theme setup during installation

### Technical
- Refactored collector to use goroutines and channels
- Added separate files for GPU, network, and battery detection
- SSH collector with parallel command execution
- Export module for multiple output formats
- Improved code organization and modularity
- Binary size: 6.2MB (includes SSH support)

## [0.1.0] - 2026-01-17

### Added
- Initial release with Go + Bubbletea foundation
- Local system info collection (OS, kernel, hostname, uptime, CPU, memory, disk, shell, terminal)
- Theme system with Lipgloss styling
- 4 built-in themes (default, minimal, dracula, nord)
- YAML configuration support
- Modular architecture for extensibility
- CLI flags for theme selection and config file
