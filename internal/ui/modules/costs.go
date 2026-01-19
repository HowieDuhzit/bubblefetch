package modules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
	"github.com/mattn/go-runewidth"
)

type ModuleCostModule struct{}

func (m *ModuleCostModule) Name() string { return "costs" }

func (m *ModuleCostModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	if info == nil || len(info.ModuleCosts) == 0 {
		return ""
	}

	costs := make([]collectors.ModuleCost, len(info.ModuleCosts))
	copy(costs, info.ModuleCosts)
	sort.SliceStable(costs, func(i, j int) bool {
		return costs[i].DurationMS > costs[j].DurationMS
	})

	maxLabel := 0
	for _, cost := range costs {
		if cost.Name == "" {
			continue
		}
		if width := runewidth.StringWidth(cost.Name); width > maxLabel {
			maxLabel = width
		}
	}

	var b strings.Builder
	b.WriteString(lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.Separator.Render("󰓅"),
		styles.Label.Render(" Module Cost"),
	))

	for _, cost := range costs {
		if cost.Name == "" {
			continue
		}
		padding := maxLabel - runewidth.StringWidth(cost.Name)
		if padding < 0 {
			padding = 0
		}
		b.WriteString("\n")
		b.WriteString("  ")
		b.WriteString(styles.Label.Render(cost.Name))
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString("  ")
		b.WriteString(styles.Value.Render(formatCost(cost.DurationMS)))
	}

	return b.String()
}

func formatCost(ms float64) string {
	switch {
	case ms >= 1000:
		return fmt.Sprintf("%.2fs", ms/1000)
	case ms >= 100:
		return fmt.Sprintf("%.0fms", ms)
	case ms >= 10:
		return fmt.Sprintf("%.1fms", ms)
	default:
		return fmt.Sprintf("%.2fms", ms)
	}
}
