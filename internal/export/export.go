package export

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
)

// ToJSON exports system info as JSON
func ToJSON(info *collectors.SystemInfo, pretty bool) (string, error) {
	var data []byte
	var err error

	if pretty {
		data, err = json.MarshalIndent(info, "", "  ")
	} else {
		data, err = json.Marshal(info)
	}

	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToYAML exports system info as YAML
func ToYAML(info *collectors.SystemInfo) (string, error) {
	data, err := yaml.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToText exports system info as plain text
func ToText(info *collectors.SystemInfo) string {
	var output string

	if info.OS != "" {
		output += fmt.Sprintf("OS: %s\n", info.OS)
	}
	if info.Kernel != "" {
		output += fmt.Sprintf("Kernel: %s\n", info.Kernel)
	}
	if info.Hostname != "" {
		output += fmt.Sprintf("Hostname: %s\n", info.Hostname)
	}
	if info.Uptime != "" {
		output += fmt.Sprintf("Uptime: %s\n", info.Uptime)
	}
	if info.CPU != "" {
		output += fmt.Sprintf("CPU: %s\n", info.CPU)
	}
	if len(info.GPU) > 0 {
		for i, gpu := range info.GPU {
			output += fmt.Sprintf("GPU %d: %s\n", i+1, gpu)
		}
	}
	if info.Memory.Total > 0 {
		output += fmt.Sprintf("Memory: %d / %d bytes\n", info.Memory.Used, info.Memory.Total)
	}
	if info.Disk.Total > 0 {
		output += fmt.Sprintf("Disk: %d / %d bytes\n", info.Disk.Used, info.Disk.Total)
	}
	if info.Shell != "" {
		output += fmt.Sprintf("Shell: %s\n", info.Shell)
	}
	if info.Terminal != "" {
		output += fmt.Sprintf("Terminal: %s\n", info.Terminal)
	}
	if info.DE != "" {
		output += fmt.Sprintf("DE: %s\n", info.DE)
	}
	if info.WM != "" {
		output += fmt.Sprintf("WM: %s\n", info.WM)
	}
	if info.LocalIP != "" {
		output += fmt.Sprintf("Local IP: %s\n", info.LocalIP)
	}
	if info.Battery.Present {
		status := "Discharging"
		if info.Battery.IsCharging {
			status = "Charging"
		}
		output += fmt.Sprintf("Battery: %.0f%% (%s)\n", info.Battery.Percentage, status)
	}

	return output
}
