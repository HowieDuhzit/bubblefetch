## Examples

This document showcases various ways to use bubblefetch.

Sample export outputs live in `examples/exports/` (HTML and SVG).

### Basic Display

```bash
# Show system info with default theme
bubblefetch

# Try different themes
bubblefetch --theme dracula
bubblefetch --theme gruvbox
bubblefetch --theme tokyo-night
bubblefetch --theme minimal
```

### Remote System Monitoring

```bash
# Check a remote server's specs
bubblefetch --remote admin@server.example.com

# Monitor a Raspberry Pi
bubblefetch --remote pi@192.168.1.100

# Check multiple servers (use a script)
for host in server1 server2 server3; do
    echo "=== $host ==="
    bubblefetch --remote $host --export text
done
```

### Export for Documentation

```bash
# Create JSON report of current system
bubblefetch --export json > system-report.json

# Generate YAML for infrastructure docs
bubblefetch --export yaml > docs/system-specs.yaml

# Create text summary
bubblefetch --export text > SYSTEM.txt

# Compact JSON (single line)
bubblefetch --export json --pretty=false > compact.json
```

### Performance Testing

```bash
# Benchmark local collection
bubblefetch --benchmark

# Benchmark remote collection
bubblefetch --remote myserver --benchmark

# Compare with other tools
echo "=== bubblefetch ==="
bubblefetch --benchmark

echo "=== neofetch ==="
time neofetch

echo "=== fastfetch ==="
time fastfetch
```

### Automation & Scripting

```bash
# Add to your shell RC file for login display
echo 'bubblefetch' >> ~/.bashrc

# Daily system report via cron
# Add to crontab: 0 9 * * * bubblefetch --export json > ~/logs/system-$(date +\%Y\%m\%d).json

# CI/CD - Capture build server specs
bubblefetch --export json > build-env.json

# Server inventory script
#!/bin/bash
for server in $(cat servers.txt); do
    bubblefetch --remote "$server" --export json > "inventory/${server}.json"
done
```

### Custom Configuration

```bash
# Create custom config
mkdir -p ~/.config/bubblefetch
cat > ~/.config/bubblefetch/config.yaml <<EOF
theme: gruvbox
modules:
  - os
  - kernel
  - uptime
  - cpu
  - gpu
  - memory
  - localip
  # - costs
EOF

# Use custom config
bubblefetch --config ~/.config/bubblefetch/config.yaml
```

### Theme Development

```bash
# Create a custom theme
cat > ~/.config/bubblefetch/themes/custom.json <<EOF
{
  "name": "custom",
  "colors": {
    "primary": "#ff6c6b",
    "secondary": "#98be65",
    "accent": "#da8548",
    "label": "#ecbe7b",
    "value": "#51afef",
    "border": "#5b6268",
    "background": "#282c34"
  },
  "ascii": "auto",
  "layout": {
    "show_ascii": true,
    "ascii_width": 30,
    "separator": " :: ",
    "padding": 2,
    "border_style": "double"
  }
}
EOF

# Test your custom theme
bubblefetch --theme custom
```

### Integration Examples

#### SSH Config Integration

```bash
# In ~/.ssh/config
Host myserver
    HostName server.example.com
    User admin
    Port 2222
    IdentityFile ~/.ssh/mykey

# Then use:
bubblefetch --remote myserver
```

#### GitHub Actions

```yaml
name: System Info
on: [push]
jobs:
  system-info:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
      - name: Install bubblefetch
        run: |
          git clone https://github.com/howieduhzit/bubblefetch.git
          cd bubblefetch && ./install.sh
      - name: Capture specs
        run: bubblefetch --export json > runner-specs.json
      - uses: actions/upload-artifact@v2
        with:
          name: system-specs
          path: runner-specs.json
```

#### Comparison Script

```bash
#!/bin/bash
# compare-systems.sh - Compare local and remote system

echo "Local System:"
bubblefetch --export text

echo ""
echo "Remote System:"
bubblefetch --remote $1 --export text

echo ""
echo "JSON Diff:"
bubblefetch --export json > /tmp/local.json
bubblefetch --remote $1 --export json > /tmp/remote.json
diff -u /tmp/local.json /tmp/remote.json || true
```

### Quick Tips

```bash
# Pipe to less for large outputs
bubblefetch --export text | less

# Save themed output (with colors)
script -c "bubblefetch" output.txt

# Quick server inventory
for i in {1..10}; do
    bubblefetch --remote server$i --export json > server$i.json &
done
wait

# Filter JSON output with jq
bubblefetch --export json | jq '.CPU, .Memory, .Disk'

# Monitor specific field
watch -n 1 'bubblefetch --export json | jq ".Memory.Used"'
```
