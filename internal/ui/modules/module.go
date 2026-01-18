package modules

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/bubblefetch/internal/collectors"
	"github.com/yourusername/bubblefetch/internal/ui/theme"
)

// Module represents a displayable system information module
type Module interface {
	Name() string
	Render(info *collectors.SystemInfo, styles theme.Styles) string
}

// pluginManager holds the global plugin manager instance
var pluginManager PluginManager

// PluginManager interface for plugin lookup
type PluginManager interface {
	GetPlugin(name string) (Module, bool)
}

// InitPlugins initializes the plugin system with the given manager
func InitPlugins(pm PluginManager) {
	pluginManager = pm
}

// Factory creates modules by name
func Factory(name string) Module {
	// Check plugins first
	if pluginManager != nil {
		if mod, ok := pluginManager.GetPlugin(name); ok {
			return mod
		}
	}

	// Fall back to built-in modules
	switch name {
	case "os":
		return &OSModule{}
	case "kernel":
		return &KernelModule{}
	case "hostname":
		return &HostnameModule{}
	case "uptime":
		return &UptimeModule{}
	case "cpu":
		return &CPUModule{}
	case "memory":
		return &MemoryModule{}
	case "disk":
		return &DiskModule{}
	case "shell":
		return &ShellModule{}
	case "terminal":
		return &TerminalModule{}
	case "de":
		return &DEModule{}
	case "wm":
		return &WMModule{}
	case "gpu":
		return &GPUModule{}
	case "network":
		return &NetworkModule{}
	case "localip":
		return &LocalIPModule{}
	case "publicip":
		return &PublicIPModule{}
	case "battery":
		return &BatteryModule{}
	default:
		return nil
	}
}

func renderField(label, value string, styles theme.Styles, separator string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.Label.Render(label),
		styles.Separator.Render(separator),
		styles.Value.Render(value),
	)
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// OSModule displays operating system information
type OSModule struct{}

func (m *OSModule) Name() string { return "os" }
func (m *OSModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return renderField("OS", info.OS, styles, ": ")
}

// KernelModule displays kernel version
type KernelModule struct{}

func (m *KernelModule) Name() string { return "kernel" }
func (m *KernelModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return renderField("Kernel", info.Kernel, styles, ": ")
}

// HostnameModule displays hostname
type HostnameModule struct{}

func (m *HostnameModule) Name() string { return "hostname" }
func (m *HostnameModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return renderField("Host", info.Hostname, styles, ": ")
}

// UptimeModule displays system uptime
type UptimeModule struct{}

func (m *UptimeModule) Name() string { return "uptime" }
func (m *UptimeModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return renderField("Uptime", info.Uptime, styles, ": ")
}

// CPUModule displays CPU information
type CPUModule struct{}

func (m *CPUModule) Name() string { return "cpu" }
func (m *CPUModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return renderField("CPU", info.CPU, styles, ": ")
}

// MemoryModule displays memory usage
type MemoryModule struct{}

func (m *MemoryModule) Name() string { return "memory" }
func (m *MemoryModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	used := formatBytes(info.Memory.Used)
	total := formatBytes(info.Memory.Total)
	value := fmt.Sprintf("%s / %s", used, total)
	return renderField("Memory", value, styles, ": ")
}

// DiskModule displays disk usage
type DiskModule struct{}

func (m *DiskModule) Name() string { return "disk" }
func (m *DiskModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	used := formatBytes(info.Disk.Used)
	total := formatBytes(info.Disk.Total)
	value := fmt.Sprintf("%s / %s", used, total)
	return renderField("Disk", value, styles, ": ")
}

// ShellModule displays shell information
type ShellModule struct{}

func (m *ShellModule) Name() string { return "shell" }
func (m *ShellModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return renderField("Shell", info.Shell, styles, ": ")
}

// TerminalModule displays terminal information
type TerminalModule struct{}

func (m *TerminalModule) Name() string { return "terminal" }
func (m *TerminalModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return renderField("Terminal", info.Terminal, styles, ": ")
}

// DEModule displays desktop environment
type DEModule struct{}

func (m *DEModule) Name() string { return "de" }
func (m *DEModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if info.DE == "" {
		return ""
	}
	return renderField("DE", info.DE, styles, ": ")
}

// WMModule displays window manager
type WMModule struct{}

func (m *WMModule) Name() string { return "wm" }
func (m *WMModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if info.WM == "" {
		return ""
	}
	return renderField("WM", info.WM, styles, ": ")
}

// GPUModule displays GPU information
type GPUModule struct{}

func (m *GPUModule) Name() string { return "gpu" }
func (m *GPUModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if len(info.GPU) == 0 {
		return ""
	}
	// Display the first GPU, or concatenate if multiple
	gpu := info.GPU[0]
	if len(info.GPU) > 1 {
		gpu = fmt.Sprintf("%s (+%d more)", gpu, len(info.GPU)-1)
	}
	return renderField("GPU", gpu, styles, ": ")
}

// NetworkModule displays network interface information
type NetworkModule struct{}

func (m *NetworkModule) Name() string { return "network" }
func (m *NetworkModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if len(info.Network) == 0 {
		return ""
	}
	// Display the first active network interface
	net := info.Network[0]
	value := fmt.Sprintf("%s (%s)", net.Interface, net.IPv4)
	return renderField("Network", value, styles, ": ")
}

// LocalIPModule displays local IP address
type LocalIPModule struct{}

func (m *LocalIPModule) Name() string { return "localip" }
func (m *LocalIPModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if info.LocalIP == "" {
		return ""
	}
	return renderField("Local IP", info.LocalIP, styles, ": ")
}

// BatteryModule displays battery information
type BatteryModule struct{}

func (m *BatteryModule) Name() string { return "battery" }
func (m *BatteryModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if !info.Battery.Present {
		return ""
	}

	status := "Discharging"
	if info.Battery.IsCharging {
		status = "Charging"
	}

	value := fmt.Sprintf("%.0f%% (%s)", info.Battery.Percentage, status)
	if info.Battery.TimeRemain != "" && !info.Battery.IsCharging {
		value = fmt.Sprintf("%.0f%% (%s remaining)", info.Battery.Percentage, info.Battery.TimeRemain)
	}

	return renderField("Battery", value, styles, ": ")
}
