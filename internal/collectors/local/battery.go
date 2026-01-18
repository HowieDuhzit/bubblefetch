package local

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yourusername/bubblefetch/internal/collectors"
)

// detectBattery attempts to detect battery information on Linux
func detectBattery() collectors.BatteryInfo {
	batteryInfo := collectors.BatteryInfo{Present: false}

	// Try to read from /sys/class/power_supply/
	powerSupplyPath := "/sys/class/power_supply"
	entries, err := os.ReadDir(powerSupplyPath)
	if err != nil {
		return batteryInfo
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "BAT") {
			batteryPath := filepath.Join(powerSupplyPath, entry.Name())
			batteryInfo.Present = true

			// Read capacity
			if data, err := os.ReadFile(filepath.Join(batteryPath, "capacity")); err == nil {
				if capacity, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
					batteryInfo.Percentage = capacity
				}
			}

			// Read status
			if data, err := os.ReadFile(filepath.Join(batteryPath, "status")); err == nil {
				status := strings.TrimSpace(string(data))
				batteryInfo.IsCharging = (status == "Charging" || status == "Full")
			}

			// Try to calculate time remaining
			if !batteryInfo.IsCharging {
				if energyNow, err := readInt(filepath.Join(batteryPath, "energy_now")); err == nil {
					if powerNow, err := readInt(filepath.Join(batteryPath, "power_now")); err == nil && powerNow > 0 {
						hoursRemain := float64(energyNow) / float64(powerNow)
						hours := int(hoursRemain)
						minutes := int((hoursRemain - float64(hours)) * 60)
						batteryInfo.TimeRemain = fmt.Sprintf("%dh %dm", hours, minutes)
					}
				}
			}

			break
		}
	}

	return batteryInfo
}

func readInt(path string) (int64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
}
