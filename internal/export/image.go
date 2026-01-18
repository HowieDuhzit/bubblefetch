package export

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/fogleman/gg"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/config"
	"github.com/howieduhzit/bubblefetch/internal/ui/modules"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

// ImageExporter handles exporting system info as images
type ImageExporter struct {
	info      *collectors.SystemInfo
	theme     *theme.Theme
	config    *config.Config
	width     int
	height    int
	fontSize  float64
	lineSpace float64
}

// NewImageExporter creates a new image exporter
func NewImageExporter(info *collectors.SystemInfo, cfg *config.Config) (*ImageExporter, error) {
	// Load theme
	thm, err := theme.Load(cfg.Theme)
	if err != nil {
		thm, _ = theme.Load("default")
	}

	return &ImageExporter{
		info:      info,
		theme:     thm,
		config:    cfg,
		width:     800,
		height:    600,
		fontSize:  14,
		lineSpace: 20,
	}, nil
}

// ToPNG exports system info as a PNG image
func (e *ImageExporter) ToPNG(outputPath string) error {
	dc := gg.NewContext(e.width, e.height)

	// Background
	bg := parseHexColor(e.theme.Colors.Background)
	dc.SetColor(bg)
	dc.Clear()

	// Load font - try multiple fallbacks
	fontPath := findFont()
	if err := dc.LoadFontFace(fontPath, e.fontSize); err != nil {
		return fmt.Errorf("failed to load font: %w", err)
	}

	// Get ASCII art
	ascii := e.theme.ASCII
	asciiLines := strings.Split(ascii, "\n")

	// Draw ASCII art (left side)
	dc.SetColor(parseHexColor(e.theme.Colors.Primary))
	x := 40.0
	y := 80.0
	for _, line := range asciiLines {
		dc.DrawString(line, x, y)
		y += e.lineSpace
	}

	// Draw system info (right side)
	infoX := 350.0
	infoY := 80.0

	// Get modules and render them
	mods := e.getModules()
	styles := e.theme.GetStyles()

	for _, mod := range mods {
		rendered := mod.Render(e.info, styles)
		if rendered == "" {
			continue
		}

		label, value, ok := splitRendered(rendered, e.theme.Layout.Separator)
		if !ok {
			continue
		}

		// Draw label
		dc.SetColor(parseHexColor(e.theme.Colors.Label))
		dc.DrawString(label+e.theme.Layout.Separator, infoX, infoY)

		// Draw value
		labelWidth, _ := dc.MeasureString(label + e.theme.Layout.Separator)
		dc.SetColor(parseHexColor(e.theme.Colors.Value))
		dc.DrawString(value, infoX+labelWidth, infoY)

		infoY += e.lineSpace + 5
	}

	// Draw border if theme has one
	if e.theme.Layout.BorderStyle != "none" {
		dc.SetColor(parseHexColor(e.theme.Colors.Border))
		dc.SetLineWidth(2)
		margin := 20.0
		dc.DrawRectangle(margin, margin, float64(e.width)-2*margin, float64(e.height)-2*margin)
		dc.Stroke()
	}

	// Save PNG
	return dc.SavePNG(outputPath)
}

// getModules returns the configured modules
func (e *ImageExporter) getModules() []modules.Module {
	mods := make([]modules.Module, 0, len(e.config.Modules))

	for _, name := range e.config.Modules {
		mod := modules.Factory(name)
		if mod != nil {
			mods = append(mods, mod)
		}
	}

	return mods
}

// parseHexColor converts a hex color string to color.Color
func parseHexColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b uint8
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)

	return color.RGBA{r, g, b, 255}
}

// findFont tries to find a suitable monospace font
func findFont() string {
	fonts := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
		"/usr/share/fonts/TTF/DejaVuSansMono.ttf",
		"/System/Library/Fonts/Monaco.ttf",
		"/System/Library/Fonts/Menlo.ttc",
		"/Library/Fonts/Courier New.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationMono-Regular.ttf",
		"/usr/share/fonts/liberation/LiberationMono-Regular.ttf",
	}

	for _, font := range fonts {
		if _, err := os.Stat(font); err == nil {
			return font
		}
	}

	// Default fallback - let gg handle it
	return "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(s string) string {
	// Simple ANSI stripping - remove escape sequences
	// This handles the lipgloss styled output
	result := ""
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(s[i])
	}

	return result
}
