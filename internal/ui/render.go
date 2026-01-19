package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/config"
	"github.com/howieduhzit/bubblefetch/internal/ui/modules"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

// Render builds the static fetch output for the given info and config.
func Render(cfg *config.Config, info *collectors.SystemInfo, err error) string {
	thm, loadErr := theme.Load(cfg.Theme)
	if loadErr != nil {
		thm, _ = theme.Load("default")
	}
	styles := thm.GetStyles()

	if err != nil {
		return styles.Value.Render("Error: " + err.Error())
	}

	var content strings.Builder

	var asciiArt string
	if thm.Layout.ShowASCII {
		asciiArt = styles.ASCII.Render(thm.ASCII)
	}

	var moduleLines []string
	for _, moduleName := range cfg.Modules {
		module := modules.Factory(moduleName)
		if module == nil {
			continue
		}
		rendered := module.Render(info, styles)
		if rendered != "" {
			moduleLines = append(moduleLines, rendered)
		}
	}

	moduleContent := strings.Join(moduleLines, "\n")

	if thm.Layout.ShowASCII && asciiArt != "" {
		content.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			asciiArt,
			strings.Repeat(" ", thm.Layout.Padding),
			moduleContent,
		))
	} else {
		content.WriteString(moduleContent)
	}

	result := styles.Container.Render(content.String())
	if thm.Layout.BorderStyle != "none" && thm.Layout.BorderStyle != "" {
		result = styles.Border.Render(result)
	}

	return result
}
