package collectors

import (
	"strings"
	"time"
)

func AddModuleCost(info *SystemInfo, name string, duration time.Duration) {
	if info == nil || name == "" {
		return
	}
	ms := float64(duration.Microseconds()) / 1000.0
	if ms < 0 {
		ms = 0
	}

	for i := range info.ModuleCosts {
		if strings.EqualFold(info.ModuleCosts[i].Name, name) {
			info.ModuleCosts[i].DurationMS = ms
			return
		}
	}

	info.ModuleCosts = append(info.ModuleCosts, ModuleCost{
		Name:       name,
		DurationMS: ms,
	})
}

func HasModuleCost(info *SystemInfo, name string) bool {
	if info == nil || name == "" {
		return false
	}
	for i := range info.ModuleCosts {
		if strings.EqualFold(info.ModuleCosts[i].Name, name) {
			return true
		}
	}
	return false
}
