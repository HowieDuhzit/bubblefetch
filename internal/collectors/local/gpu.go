package local

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// detectGPU attempts to detect GPU information
func detectGPU() []string {
	var gpus []string

	// Try reading /sys/class/drm FIRST (much faster than lspci)
	drmPath := "/sys/class/drm"
	if entries, err := os.ReadDir(drmPath); err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "card") && !strings.Contains(entry.Name(), "-") {
				devicePath := filepath.Join(drmPath, entry.Name(), "device")

				// Try to read uevent file which has GPU info
				if uevent, err := os.ReadFile(filepath.Join(devicePath, "uevent")); err == nil {
					lines := strings.Split(string(uevent), "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, "PCI_ID=") {
							// Found PCI ID, now try to get a friendly name
							gpus = append(gpus, parseGPUFromSys(devicePath))
							break
						}
					}
				}
			}
		}
	}

	// If we found GPUs via /sys, return them
	if len(gpus) > 0 {
		return gpus
	}

	// Fall back to lspci with timeout (only if /sys failed)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, "lspci")
	if output, err := cmd.Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			lowerLine := strings.ToLower(line)
			if strings.Contains(lowerLine, "vga") || strings.Contains(lowerLine, "3d") || strings.Contains(lowerLine, "display") {
				// Extract GPU name (after the colon)
				parts := strings.SplitN(line, ":", 2)
				if len(parts) > 1 {
					gpu := strings.TrimSpace(parts[1])
					// Clean up common prefixes
					gpu = strings.TrimPrefix(gpu, "VGA compatible controller: ")
					gpu = strings.TrimPrefix(gpu, "3D controller: ")
					gpu = strings.TrimPrefix(gpu, "Display controller: ")
					gpus = append(gpus, gpu)
				}
			}
		}
	}

	// Final fallback - return empty to avoid showing "Unknown"
	return gpus
}

func parseGPUFromSys(devicePath string) string {
	// Read PCI ID from uevent
	if uevent, err := os.ReadFile(filepath.Join(devicePath, "uevent")); err == nil {
		lines := strings.Split(string(uevent), "\n")
		var vendor, device string

		for _, line := range lines {
			if strings.HasPrefix(line, "PCI_ID=") {
				parts := strings.Split(strings.TrimPrefix(line, "PCI_ID="), ":")
				if len(parts) == 2 {
					vendor = parts[0]
					device = parts[1]
				}
			}
		}

		if vendor != "" && device != "" {
			// Map common vendor IDs to names
			vendorName := vendor
			switch vendor {
			case "10DE":
				vendorName = "NVIDIA"
			case "1002":
				vendorName = "AMD"
			case "8086":
				vendorName = "Intel"
			}
			return vendorName + " GPU (" + device + ")"
		}
	}

	return "GPU"
}

func readFirstLine(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}
