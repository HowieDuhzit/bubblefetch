package export

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>bubblefetch - System Information</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            background: {{.Colors.Background}};
            color: {{.Colors.Value}};
            font-family: 'Courier New', 'Consolas', 'Monaco', monospace;
            padding: 40px 20px;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .container {
            display: flex;
            max-width: 1000px;
            width: 100%;
            {{if ne .BorderStyle "none"}}
            border: 2px solid {{.Colors.Border}};
            {{end}}
            border-radius: 10px;
            padding: 40px;
            gap: 40px;
        }

        .ascii {
            color: {{.Colors.Primary}};
            white-space: pre;
            flex: 0 0 auto;
            line-height: 1.4;
            font-size: 14px;
        }

        .info {
            flex: 1;
            min-width: 0;
        }

        .field {
            margin-bottom: 8px;
            display: flex;
            gap: 4px;
        }

        .label {
            color: {{.Colors.Label}};
            font-weight: bold;
            white-space: nowrap;
        }

        .separator {
            color: {{.Colors.Label}};
        }

        .value {
            color: {{.Colors.Value}};
            word-break: break-word;
        }

        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid {{.Colors.Border}};
            color: {{.Colors.Secondary}};
            font-size: 12px;
            text-align: center;
        }

        @media (max-width: 768px) {
            .container {
                flex-direction: column;
                gap: 20px;
            }

            .ascii {
                font-size: 10px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="ascii">{{.ASCII}}</div>
        <div class="info">
            {{range .Fields}}
            <div class="field">
                <span class="label">{{.Label}}</span>
                <span class="separator">{{$.Separator}}</span>
                <span class="value">{{.Value}}</span>
            </div>
            {{end}}
            <div class="footer">
                Generated with bubblefetch v{{.Version}} | Theme: {{.ThemeName}}
            </div>
        </div>
    </div>
</body>
</html>`

type htmlField struct {
	Label string
	Value string
}

// ToHTML exports system info as an HTML file
func (e *ImageExporter) ToHTML(outputPath string) error {
	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	data := e.prepareHTMLData()

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	return nil
}

func (e *ImageExporter) prepareHTMLData() map[string]interface{} {
	ascii := e.theme.ASCII

	// Get module data
	fields := []htmlField{}
	mods := e.getModules()
	styles := e.theme.GetStyles()

	for _, mod := range mods {
		rendered := mod.Render(e.info, styles)
		if rendered == "" {
			continue
		}

		parts := strings.SplitN(rendered, e.theme.Layout.Separator, 2)
		if len(parts) != 2 {
			continue
		}

		label := stripANSI(parts[0])
		value := stripANSI(parts[1])

		fields = append(fields, htmlField{
			Label: label,
			Value: value,
		})
	}

	return map[string]interface{}{
		"Colors":     e.theme.Colors,
		"BorderStyle": e.theme.Layout.BorderStyle,
		"ASCII":      ascii,
		"Fields":     fields,
		"Separator":  e.theme.Layout.Separator,
		"Version":    "0.3.0",
		"ThemeName":  e.config.Theme,
	}
}
