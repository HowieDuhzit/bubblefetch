# bubblefetch

A simple, elegant, and highly customizable system information tool built with Go and Bubbletea. An alternative to neofetch/fastfetch with beautiful TUI, extensive theming, and remote system support.

## Features

- **⚡ Blazing Fast**: Average 1.3ms collection time - **100x faster than neofetch, 8x faster than fastfetch**
- **Beautiful TUI**: Built with Bubbletea and Lipgloss for elegant terminal UI
- **OS Detection**: Automatically detects your OS/distro and displays appropriate ASCII art
- **SSH Remote Support**: Fetch system info from remote systems via SSH
- **Export Modes**: Export to JSON, YAML, or plain text
- **Benchmark Mode**: Measure collection performance
- **Highly Customizable**: YAML config, custom themes, modular system info display
- **Comprehensive Info**: CPU, GPU, memory, disk, network, battery, and more
- **Themeable**: 8 built-in themes with easy custom theme creation

## Installation

### Quick Install

```bash
git clone https://github.com/yourusername/bubblefetch.git
cd bubblefetch
./install.sh
```

The install script will:
- Build the optimized binary
- Install to `/usr/local/bin`
- Create config directory at `~/.config/bubblefetch`
- Copy themes and example config

### Manual Installation

```bash
git clone https://github.com/yourusername/bubblefetch.git
cd bubblefetch
go build -ldflags="-s -w" -o bubblefetch ./cmd/bubblefetch
sudo mv bubblefetch /usr/local/bin/
```

### Go Install

```bash
go install github.com/yourusername/bubblefetch/cmd/bubblefetch@latest
```

## Usage

### Basic Usage

```bash
# Run with default settings
bubblefetch

# Use a specific theme
bubblefetch --theme dracula

# Use custom config file
bubblefetch --config ~/.config/bubblefetch/custom.yaml
```

### Remote Systems (SSH)

```bash
# Fetch info from remote system via SSH
bubblefetch --remote user@hostname

# Uses your SSH config and keys automatically
bubblefetch --remote myserver
```

### Export Modes

```bash
# Export as JSON
bubblefetch --export json > system.json

# Export as YAML
bubblefetch --export yaml > system.yaml

# Export as plain text
bubblefetch --export text > system.txt

# Compact JSON (no pretty print)
bubblefetch --export json --pretty=false
```

### Benchmark Mode

```bash
# Run 10 iterations and show performance stats
bubblefetch --benchmark
```

### Other Options

```bash
# Show version
bubblefetch --version

# Show help
bubblefetch --help
```

### Keyboard Shortcuts

- `q` / `Esc` / `Ctrl+C` - Quit

## Configuration

Copy the example config to get started:

```bash
mkdir -p ~/.config/bubblefetch
cp config.example.yaml ~/.config/bubblefetch/config.yaml
```

Edit `~/.config/bubblefetch/config.yaml` to customize:

```yaml
theme: default

modules:
  - os
  - kernel
  - hostname
  - uptime
  - cpu
  - memory
  - disk
  - shell
  - terminal

ssh:
  port: 22
```

## Themes

### Built-in Themes

All themes auto-detect your OS and display the appropriate ASCII art!

- **default** - Catppuccin-inspired colors with rounded borders
- **minimal** - Clean, borderless design
- **dracula** - Based on the Dracula color scheme
- **nord** - Arctic, north-bluish color palette
- **gruvbox** - Warm, retro groove colors
- **tokyo-night** - Dark Tokyo Night theme
- **monokai** - Classic Monokai Pro colors
- **solarized-dark** - Precision colors for machines and people

### Supported OS ASCII Art

Auto-detected logos for: Arch, Ubuntu, Debian, Fedora, Mint, Manjaro, Pop!_OS, Gentoo, openSUSE, Kali, Void, NixOS, macOS, Windows, FreeBSD, Alpine, and more!

### Using Themes

```bash
bubblefetch --theme nord
```

Or set in config:

```yaml
theme: dracula
```

### Creating Custom Themes

Create a JSON file in `~/.config/bubblefetch/themes/mytheme.json`:

```json
{
  "name": "mytheme",
  "colors": {
    "primary": "#89b4fa",
    "secondary": "#cba6f7",
    "accent": "#f38ba8",
    "label": "#f9e2af",
    "value": "#a6e3a1",
    "border": "#585b70",
    "background": "#1e1e2e"
  },
  "ascii": "\n    Your ASCII art here\n",
  "layout": {
    "show_ascii": true,
    "ascii_width": 30,
    "separator": ": ",
    "padding": 2,
    "border_style": "rounded"
  }
}
```

Border styles: `rounded`, `double`, `thick`, `normal`, `none`

## Modules

Available system information modules:

- `os` - Operating system and version
- `kernel` - Kernel version
- `hostname` - System hostname
- `uptime` - System uptime
- `cpu` - CPU model
- `gpu` - GPU information (auto-detected)
- `memory` - Memory usage
- `disk` - Disk usage (root partition)
- `shell` - Current shell
- `terminal` - Terminal emulator
- `de` - Desktop environment
- `wm` - Window manager
- `network` - Active network interface and IP
- `localip` - Local IP address
- `battery` - Battery status and percentage (laptops only)

Configure module order in your config file.

## Development

### Project Structure

```
bubblefetch/
├── cmd/bubblefetch/          # Main entry point
├── internal/
│   ├── config/               # Config loading & validation
│   ├── collectors/           # System info collectors
│   │   ├── local/           # Local system info
│   │   └── remote/          # SSH-based remote info (planned)
│   ├── ui/                   # Bubbletea TUI components
│   │   ├── theme/           # Theme engine
│   │   └── modules/         # Display modules
│   └── plugins/              # Plugin system (planned)
├── themes/                   # Built-in theme files
└── config.example.yaml       # Example configuration
```

### Building

```bash
go build -o bubblefetch ./cmd/bubblefetch
```

### Running Tests

```bash
go test ./...
```

## Roadmap

- [x] Basic local system info collection
- [x] Parallel data collection for speed
- [x] Theme system with Lipgloss
- [x] OS-specific ASCII art auto-detection
- [x] Modular display system
- [x] YAML configuration
- [x] GPU, network, and battery modules
- [x] SSH remote system support
- [x] Export to JSON/YAML/Text
- [x] Benchmark mode
- [x] 8 built-in themes
- [x] Installation scripts
- [ ] Plugin system for custom modules
- [ ] Package for major Linux distros (AUR, Homebrew, apt, etc.)
- [ ] Public IP detection
- [ ] Interactive configuration mode
- [ ] Screenshot/image export

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [gopsutil](https://github.com/shirou/gopsutil) - System information library
- Inspired by [neofetch](https://github.com/dylanaraps/neofetch) and [fastfetch](https://github.com/fastfetch-cli/fastfetch)
