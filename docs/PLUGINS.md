# bubblefetch Plugin System

The bubblefetch plugin system allows you to create custom display modules using Go plugins (.so files). Plugins extend bubblefetch with new modules that can display custom information in the TUI.

## Platform Support

**Supported Platforms:**
- ✅ Linux
- ✅ macOS
- ✅ FreeBSD

**Not Supported:**
- ❌ Windows (Go plugin system limitation)

For Windows users, we recommend contributing new modules directly to the bubblefetch core instead.

## Quick Start

### 1. Create a Plugin

Create a file `myplugin.go`:

```go
package main

import (
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

// ModuleName must be a package-level variable named exactly "ModuleName"
var ModuleName = "myplugin"

// Render must have this exact signature
func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	// Access system info
	hostname := info.Hostname

	// Format output with theme styles
	label := styles.Label.Render("Custom")
	separator := styles.Separator.Render(": ")
	value := styles.Value.Render("Hello from " + hostname)

	return label + separator + value
}
```

### 2. Build the Plugin

```bash
go build -buildmode=plugin -o myplugin.so myplugin.go
```

### 3. Install the Plugin

```bash
mkdir -p ~/.config/bubblefetch/plugins
cp myplugin.so ~/.config/bubblefetch/plugins/
```

### 4. Use the Plugin

Add to your `~/.config/bubblefetch/config.yaml`:

```yaml
modules:
  - os
  - myplugin  # Your custom plugin!
  - cpu
```

Run bubblefetch and your plugin will appear in the TUI!

## Plugin Registry

The GitHub Pages site reads `plugins/manifest.json` to list plugins in the
browser. Downloads should point to GitHub Release assets so the repo does not
store binary `.so` files.

Manifest example:

```json
{
  "id": "hello",
  "name": "Hello",
  "version": "0.1.0",
  "downloads": [
    {
      "label": "Linux (amd64)",
      "url": "https://github.com/howieduhzit/bubblefetch/releases/download/plugins-v0.1.0/hello_linux_amd64.so"
    }
  ],
  "source": "https://github.com/howieduhzit/bubblefetch/blob/main/plugins/examples/hello.go"
}
```

## Building a release asset

Go plugins are OS/arch specific. Build on the target platform and upload the
resulting `.so` to a GitHub Release.

```bash
scripts/build-plugin.sh hello plugins/examples/hello.go
```

Then update `plugins/manifest.json` with the release asset URL.

## Plugin API

### Required Exports

Every plugin must export exactly these two symbols:

#### 1. ModuleName (variable)

```go
var ModuleName = "pluginname"
```

- **Type:** `string` (must be a pointer to string when loaded)
- **Purpose:** Identifies the plugin module
- **Used in:** config.yaml `modules` list
- **Naming:** Use lowercase, alphanumeric characters

#### 2. Render (function)

```go
func Render(info *collectors.SystemInfo, styles theme.Styles) string
```

- **Parameters:**
  - `info *collectors.SystemInfo` - All collected system information
  - `styles theme.Styles` - Current theme's style definitions
- **Returns:** `string` - The formatted output to display
- **Purpose:** Renders the module's output for the TUI

### Available System Information

The `collectors.SystemInfo` struct provides access to:

```go
type SystemInfo struct {
	// Basic info
	OS          string
	Kernel      string
	Hostname    string
	Uptime      string

	// Hardware
	CPU         string
	GPU         []string
	Memory      MemoryInfo
	Disk        DiskInfo

	// Environment
	Shell       string
	Terminal    string
	DE          string  // Desktop Environment
	WM          string  // Window Manager

	// Network
	Network     []NetworkInfo
	LocalIP     string
	PublicIP    string

	// Power
	Battery     BatteryInfo
}
```

See `internal/collectors/types.go` for complete struct definitions.

### Theme Styles

The `theme.Styles` provides these styled renderers:

```go
styles.Label.Render("Label")        // Renders label text
styles.Value.Render("Value")        // Renders value text
styles.Separator.Render(": ")       // Renders separator
styles.Primary.Render("Text")       // Primary color
styles.Secondary.Render("Text")     // Secondary color
styles.Accent.Render("Text")        // Accent color
```

Using theme styles ensures your plugin respects the user's theme choice.

## Example Plugins

### Simple Greeting

```go
package main

import (
	"fmt"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

var ModuleName = "hello"

func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	greeting := fmt.Sprintf("Hello from %s!", info.Hostname)

	label := styles.Label.Render("Greeting")
	separator := styles.Separator.Render(": ")
	value := styles.Value.Render(greeting)

	return label + separator + value
}
```

### CPU Core Count

```go
package main

import (
	"fmt"
	"runtime"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

var ModuleName = "cores"

func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	cores := runtime.NumCPU()

	label := styles.Label.Render("CPU Cores")
	separator := styles.Separator.Render(": ")
	value := styles.Value.Render(fmt.Sprintf("%d", cores))

	return label + separator + value
}
```

### Custom System Info

```go
package main

import (
	"fmt"
	"os"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

var ModuleName = "myinfo"

func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	// Read custom info from environment or file
	customValue := os.Getenv("MY_CUSTOM_VAR")
	if customValue == "" {
		customValue = "Not set"
	}

	label := styles.Label.Render("My Info")
	separator := styles.Separator.Render(": ")
	value := styles.Value.Render(customValue)

	return label + separator + value
}
```

## Building Plugins

### Manual Build

```bash
go build -buildmode=plugin -o myplugin.so myplugin.go
```

### Using Makefile

Add to `Makefile`:

```makefile
plugin-myplugin:
	go build -buildmode=plugin -o plugins/myplugin.so plugins/myplugin.go
```

Then:

```bash
make plugin-myplugin
```

### Dependencies

Plugins can import external packages:

```bash
# In your plugin directory
go mod init myplugin
go get github.com/some/package

go build -buildmode=plugin -o myplugin.so .
```

**Important:** Plugins must be built with the **same Go version** and **same dependency versions** as bubblefetch to avoid runtime errors.

## Configuration

### Plugin Directory

Default: `~/.config/bubblefetch/plugins/`

Custom directory in `config.yaml`:

```yaml
plugin_dir: /path/to/my/plugins
```

### Loading Plugins

Plugins are automatically loaded from the plugin directory when bubblefetch starts. Any `.so` file in the directory will be loaded.

### Using Plugins

Add plugin names to the `modules` list:

```yaml
modules:
  - os
  - kernel
  - myplugin1
  - cpu
  - myplugin2
  - memory
```

Order matters - modules display in the order listed.

## Error Handling

### Plugin Load Failures

If a plugin fails to load, bubblefetch will:
1. Print a warning to stderr
2. Continue loading other plugins
3. Skip the failed plugin module

Common errors:

```
Warning: failed to load plugin myplugin.so: plugin missing ModuleName
Warning: failed to load plugin myplugin.so: Render is not the correct function type
```

### Debugging Plugins

Run with stderr visible:

```bash
bubblefetch 2>&1 | grep -i plugin
```

Check if plugin loaded:

```bash
ls -la ~/.config/bubblefetch/plugins/
```

Verify plugin symbols:

```bash
nm -D myplugin.so | grep -E "ModuleName|Render"
```

## Best Practices

### 1. Keep Plugins Simple

Plugins should be lightweight and fast. Heavy computation can slow down bubblefetch startup.

```go
// Good - fast and simple
func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return styles.Label.Render("User") + styles.Separator.Render(": ") +
	       styles.Value.Render(os.Getenv("USER"))
}

// Bad - slow external API call
func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	weather := fetchWeatherAPI() // Blocks startup!
	return formatWeather(weather)
}
```

### 2. Handle Errors Gracefully

Return empty string or fallback text on errors:

```go
func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	data, err := readCustomData()
	if err != nil {
		// Don't panic - return empty or fallback
		return ""
	}
	return formatData(data, styles)
}
```

### 3. Use Theme Styles

Always use the provided styles instead of hardcoded colors:

```go
// Good - respects user's theme
label := styles.Label.Render("Custom")
value := styles.Value.Render(data)

// Bad - hardcoded colors break themes
label := "\033[1;34mCustom\033[0m"
```

### 4. Return Empty for Missing Data

If your plugin has no data to show, return empty string:

```go
func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	data := getOptionalData()
	if data == "" {
		return "" // Module won't display
	}
	return formatData(data, styles)
}
```

### 5. Match Existing Module Format

Follow the same format as built-in modules:

```
Label: Value
```

Use the helper pattern:

```go
label := styles.Label.Render("Label")
separator := styles.Separator.Render(": ")
value := styles.Value.Render("Value")
return label + separator + value
```

## Limitations

### Plugins are Display-Only

Plugins cannot:
- Add new data to `SystemInfo` struct
- Modify existing system info
- Affect JSON/YAML/text export modes
- Change theme styles
- Modify other modules

Plugins can only:
- Read existing `SystemInfo` data
- Format and display custom output for the TUI

### No Inter-Plugin Communication

Plugins cannot communicate with each other or share state.

### Platform Restrictions

Go plugins only work on Linux, macOS, and FreeBSD. Windows is not supported.

## Advanced Examples

### External API Call (Cached)

```go
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

var ModuleName = "weather"

var (
	weatherCache     string
	weatherCacheTime time.Time
	weatherMutex     sync.Mutex
)

func Render(info *collectors.SystemInfo, styles theme.Styles) string {
	weatherMutex.Lock()
	defer weatherMutex.Unlock()

	// Cache for 30 minutes
	if time.Since(weatherCacheTime) > 30*time.Minute {
		weatherCache = fetchWeather()
		weatherCacheTime = time.Now()
	}

	if weatherCache == "" {
		return ""
	}

	label := styles.Label.Render("Weather")
	separator := styles.Separator.Render(": ")
	value := styles.Value.Render(weatherCache)
	return label + separator + value
}

func fetchWeather() string {
	// Quick timeout to avoid blocking
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://wttr.in?format=j1")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Parse and return weather
	// (implementation details...)
	return "Sunny, 72°F"
}
```

## Troubleshooting

### Plugin Not Loading

1. **Check file extension**: Must be `.so`
2. **Check permissions**: `chmod +x myplugin.so`
3. **Check location**: Must be in `~/.config/bubblefetch/plugins/`
4. **Check symbols**: Run `nm -D myplugin.so | grep -E "ModuleName|Render"`

### Plugin Not Appearing

1. **Check config**: Module name must be in `modules` list
2. **Check return value**: Plugin must return non-empty string
3. **Check spelling**: Module name must match `ModuleName` variable exactly
4. **Check stderr**: Look for load warnings

### Symbol Errors

```
Error: ModuleName is not *string
```

Make sure `ModuleName` is a package-level **variable**, not a constant:

```go
// Correct
var ModuleName = "myplugin"

// Wrong
const ModuleName = "myplugin"
```

### Type Mismatch Errors

```
Error: Render is not the correct function type
```

Signature must match exactly:

```go
func Render(info *collectors.SystemInfo, styles theme.Styles) string
```

## Contributing Plugins

If you create a useful plugin, consider contributing it:

1. **Share in Discussions**: Post in GitHub Discussions
2. **Create Example**: Add to `plugins/examples/`
3. **Submit PR**: Contribute to bubblefetch core

Popular or widely useful plugins may be integrated into bubblefetch as built-in modules.

## Resources

- **Example Plugins**: `plugins/examples/`
- **SystemInfo Struct**: `internal/collectors/types.go`
- **Theme Styles**: `internal/ui/theme/theme.go`
- **Built-in Modules**: `internal/ui/modules/`
- **Go Plugins**: https://pkg.go.dev/plugin

## Summary

1. Create plugin with `ModuleName` var and `Render` function
2. Build with `go build -buildmode=plugin`
3. Copy to `~/.config/bubblefetch/plugins/`
4. Add module name to config
5. Enjoy your custom module!

Happy plugin development! 🔌
