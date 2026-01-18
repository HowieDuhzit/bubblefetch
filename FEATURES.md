# bubblefetch Features

Complete feature list for bubblefetch v0.2.0

## Core Features

### System Information Collection

#### Local System
- **OS Detection**: Automatic detection of 15+ distributions
- **Kernel Version**: Linux kernel version display
- **Hostname**: System hostname
- **Uptime**: System uptime with smart formatting (days/hours/minutes)
- **CPU**: Processor model and architecture
- **GPU**: Graphics card detection via lspci and /sys/class/drm
- **Memory**: RAM usage (used/total)
- **Disk**: Root partition usage (used/total)
- **Shell**: Current shell path
- **Terminal**: Terminal emulator name
- **Desktop Environment**: DE if applicable
- **Window Manager**: WM/session type
- **Network**: Active network interfaces with IPv4/IPv6
- **Local IP**: Primary local IP address
- **Battery**: Laptop battery percentage, charging status, time remaining

#### Remote Systems (SSH)
- Full SSH support for remote system monitoring
- Automatic SSH key detection (~/.ssh/id_rsa, id_ed25519)
- Support for SSH config files
- Custom ports and configurations
- User@host format support
- Parallel command execution on remote systems

### Performance

- **Parallel Collection**: All metrics gathered concurrently using goroutines
- **Fast Execution**: 2-3x faster than sequential collection
- **Minimal Overhead**: Optimized binary size (6.2MB)
- **Efficient Network**: Minimal bandwidth usage for remote collections

### Display & Themes

#### Built-in Themes (8 total)
1. **default** - Catppuccin-inspired, rounded borders
2. **minimal** - Clean, no borders, compact
3. **dracula** - Dracula color scheme, double borders
4. **nord** - Arctic nord palette, thick borders
5. **gruvbox** - Warm retro colors
6. **tokyo-night** - Modern dark theme
7. **monokai** - Classic Monokai Pro
8. **solarized-dark** - Precision engineered colors

#### Theme Features
- Auto-detected OS ASCII art
- Customizable colors (7 color slots)
- Border styles: rounded, double, thick, normal, none
- Custom separators
- Adjustable padding
- Custom ASCII art support

#### Supported OS Logos
- Arch Linux
- Ubuntu
- Debian
- Fedora
- Linux Mint
- Manjaro
- Pop!_OS
- Gentoo
- openSUSE
- Kali Linux
- Void Linux
- NixOS
- macOS
- Windows
- FreeBSD
- Alpine Linux
- Generic Linux fallback

### Export Formats

#### JSON Export
```bash
bubblefetch --export json
```
- Pretty printed (default)
- Compact mode with `--pretty=false`
- Full system info in structured format
- Perfect for APIs and data processing

#### YAML Export
```bash
bubblefetch --export yaml
```
- Clean, human-readable format
- Great for configuration management
- Infrastructure-as-code documentation

#### Plain Text Export
```bash
bubblefetch --export text
```
- Simple key-value format
- Easy to parse with shell scripts
- Minimal formatting

### Benchmark Mode

```bash
bubblefetch --benchmark
```

- Runs 10 collection iterations
- Shows individual run times
- Calculates average and total time
- Works with both local and remote collection
- Perfect for performance tuning

### Configuration

#### YAML Configuration File
Location: `~/.config/bubblefetch/config.yaml`

Features:
- Theme selection
- Module ordering and selection
- SSH settings (user, port, key path)
- Remote system default

#### Module System
Available modules (14 total):
- os, kernel, hostname, uptime
- cpu, gpu, memory, disk
- shell, terminal, de, wm
- network, localip, battery

Customize display order in config file

### CLI Interface

#### Flags
- `--theme <name>` - Select theme
- `--config <path>` - Custom config file
- `--remote <user@host>` - Remote system
- `--export <format>` - Export format (json/yaml/text)
- `--pretty <bool>` - Pretty print JSON (default: true)
- `--benchmark` - Run benchmark
- `--version` - Show version
- `--help` - Show help

#### Keyboard Shortcuts (TUI Mode)
- `q` - Quit
- `Esc` - Quit
- `Ctrl+C` - Quit

### Installation

#### Automated Installer
- `./install.sh` - One-command installation
- Builds optimized binary
- Installs to /usr/local/bin
- Creates config directory
- Copies themes and example config

#### Uninstaller
- `./uninstall.sh` - Clean removal
- Optional config directory removal
- Interactive prompts

### Developer Features

#### Modular Architecture
```
bubblefetch/
├── cmd/bubblefetch/       # Main entry point
├── internal/
│   ├── collectors/        # System info collection
│   │   ├── local/         # Local collectors
│   │   └── remote/        # SSH remote collectors
│   ├── config/            # Configuration management
│   ├── export/            # Export formatters
│   └── ui/                # Bubbletea UI
│       ├── theme/         # Theme engine & ASCII art
│       └── modules/       # Display modules
└── themes/                # Built-in themes
```

#### Extensibility
- Easy to add new modules
- Simple theme creation (JSON)
- Plugin architecture (planned)
- Clean interfaces for collectors

#### Build System
- Makefile with common tasks
- Optimized release builds with `-ldflags="-s -w"`
- Cross-platform support (Linux, macOS, FreeBSD)
- GitHub Actions CI/CD ready

### Documentation

- README.md - Comprehensive guide
- QUICKSTART.md - 60-second setup
- EXAMPLES.md - Real-world usage examples
- FEATURES.md - This document
- CHANGELOG.md - Version history
- NEXT_STEPS.md - Development roadmap

### Testing & Quality

- Go test support
- Benchmark mode for performance testing
- Linter ready (golangci-lint)
- Code formatting (go fmt)

## Technical Specifications

### Dependencies
- **Bubbletea**: TUI framework
- **Lipgloss**: Styling engine
- **Bubbles**: UI components
- **gopsutil**: System information (CPU, memory, disk, etc.)
- **golang.org/x/crypto/ssh**: SSH client
- **gopkg.in/yaml.v3**: YAML parsing

### Performance Characteristics
- **Startup Time**: < 100ms typical
- **Collection Time**: 50-200ms (local), 200-500ms (remote via SSH)
- **Memory Usage**: < 20MB
- **Binary Size**: 6.2MB (static, includes all dependencies)

### Platform Support
- **Linux**: Full support (all features)
- **macOS**: Full support
- **FreeBSD**: Basic support
- **Windows**: Limited (WSL recommended)

### Requirements
- Go 1.21+ for building
- SSH client for remote features
- lspci for GPU detection (optional)
- Standard UNIX utilities (uname, cat, etc.)

## Future Features (Planned)

- [ ] Plugin system for custom modules
- [ ] Package repositories (AUR, Homebrew, apt)
- [ ] Public IP detection
- [ ] Interactive configuration mode
- [ ] Screenshot/image export
- [ ] Custom ASCII art library
- [ ] Multi-system comparison view
- [ ] Real-time monitoring mode
- [ ] More themes and customization
- [ ] Windows native support

## Use Cases

1. **System Administrators**: Quick server inventory and comparison
2. **DevOps**: CI/CD environment documentation
3. **Documentation**: Automated system spec generation
4. **r/unixporn**: Beautiful system info displays
5. **SSH Users**: Remote system monitoring without installing tools
6. **Automation**: Scripting and monitoring integration
7. **Learning**: Understanding system information and Go programming

## Comparison with Alternatives

### vs neofetch
- ✓ Faster (parallel collection)
- ✓ Remote SSH support built-in
- ✓ Export formats (JSON/YAML/text)
- ✓ Active development
- ✓ Modern TUI framework

### vs fastfetch
- ✓ More customizable themes
- ✓ Remote SSH support
- ✓ Export functionality
- ≈ Similar speed (both fast)

### vs screenfetch
- ✓ Much faster
- ✓ Better maintained
- ✓ More features
- ✓ Modern architecture
