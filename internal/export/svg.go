package export

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

const svgTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="{{.Width}}" height="{{.Height}}" xmlns="http://www.w3.org/2000/svg">
    <style>
        .background { fill: {{.Colors.Background}}; }
        .border { fill: none; stroke: {{.Colors.Border}}; stroke-width: 2; }
        .ascii { fill: {{.Colors.Primary}}; font-family: monospace; font-size: 14px; }
        .label { fill: {{.Colors.Label}}; font-family: monospace; font-size: 14px; font-weight: bold; }
        .value { fill: {{.Colors.Value}}; font-family: monospace; font-size: 14px; }
    </style>

    <!-- Background -->
    <rect class="background" width="{{.Width}}" height="{{.Height}}" />

    {{if ne .BorderStyle "none"}}
    <!-- Border -->
    <rect class="border" x="20" y="20" width="{{sub .Width 40}}" height="{{sub .Height 40}}" rx="10" />
    {{end}}

    <!-- ASCII Art -->
    {{range $i, $line := .ASCIILines}}
    <text class="ascii" x="40" y="{{add 80 (mul $i 20)}}">{{escape $line}}</text>
    {{end}}

    <!-- System Info -->
    {{range $i, $field := .Fields}}
    <text class="label" x="350" y="{{add 80 (mul $i 25)}}">{{escape $field.Label}}{{$.Separator}}</text>
    <text class="value" x="{{add 350 $field.LabelWidth}}" y="{{add 80 (mul $i 25)}}">{{escape $field.Value}}</text>
    {{end}}
</svg>`

// ToSVG exports system info as an SVG image
func (e *ImageExporter) ToSVG(outputPath string) error {
	tmpl, err := template.New("svg").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"escape": func(s string) string {
			s = strings.ReplaceAll(s, "&", "&amp;")
			s = strings.ReplaceAll(s, "<", "&lt;")
			s = strings.ReplaceAll(s, ">", "&gt;")
			return s
		},
	}).Parse(svgTemplate)

	if err != nil {
		return fmt.Errorf("failed to parse SVG template: %w", err)
	}

	data := e.prepareSVGData()

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute SVG template: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write SVG file: %w", err)
	}

	return nil
}

type svgField struct {
	Label      string
	Value      string
	LabelWidth int
}

func (e *ImageExporter) prepareSVGData() map[string]interface{} {
	ascii := e.theme.ASCII
	rawAsciiLines := strings.Split(ascii, "\n")
	asciiLines := make([]string, 0, len(rawAsciiLines))
	for _, line := range rawAsciiLines {
		asciiLines = append(asciiLines, sanitizeText(stripANSI(line)))
	}

	// Get module data
	fields := []svgField{}
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

		label = sanitizeText(label)
		value = sanitizeText(value)

		// Estimate label width for SVG positioning (rough approximation)
		labelWidth := len(label+e.theme.Layout.Separator) * 8

		fields = append(fields, svgField{
			Label:      label,
			Value:      value,
			LabelWidth: labelWidth,
		})
	}

	return map[string]interface{}{
		"Width":       e.width,
		"Height":      e.height,
		"Colors":      e.theme.Colors,
		"BorderStyle": e.theme.Layout.BorderStyle,
		"ASCIILines":  asciiLines,
		"Fields":      fields,
		"Separator":   e.theme.Layout.Separator,
	}
}
