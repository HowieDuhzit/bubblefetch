package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/modules"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

type externalModule struct {
	name    string
	path    string
	timeout time.Duration
}

type externalOutput struct {
	Label     string   `json:"label"`
	Value     string   `json:"value"`
	Icon      string   `json:"icon"`
	Separator string   `json:"separator"`
	Lines     []string `json:"lines"`
	Raw       string   `json:"raw"`
	Text      string   `json:"text"`
}

func (m *externalModule) Name() string {
	return m.name
}

func (m *externalModule) Render(_ *collectors.SystemInfo, styles theme.Styles) string {
	output, err := m.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: external module %s failed: %v\n", m.name, err)
		return ""
	}

	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return ""
	}

	if strings.HasPrefix(trimmed, "{") {
		var payload externalOutput
		if err := json.Unmarshal([]byte(trimmed), &payload); err == nil {
			return formatExternal(payload, styles)
		}
	}

	return formatExternal(externalOutput{Raw: trimmed}, styles)
}

func (m *externalModule) run() (string, error) {
	timeout := m.timeout
	if timeout <= 0 {
		timeout = 250 * time.Millisecond
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, m.path)
	cwd, err := os.Getwd()
	if err == nil && cwd != "" {
		cmd.Dir = cwd
	}
	cmd.Env = append(os.Environ(),
		"BUBBLEFETCH_FORMAT=json",
		"BUBBLEFETCH_MODULE="+m.name,
		"BUBBLEFETCH_CWD="+cmd.Dir,
	)

	output, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout after %s", timeout)
	}
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func formatExternal(payload externalOutput, styles theme.Styles) string {
	separator := payload.Separator
	if separator == "" {
		separator = ": "
	}

	if payload.Label != "" {
		if payload.Icon != "" {
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				styles.Separator.Render(payload.Icon),
				styles.Label.Render(" "+payload.Label),
				styles.Separator.Render(separator),
				styles.Value.Render(payload.Value),
			)
		}
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.Label.Render(payload.Label),
			styles.Separator.Render(separator),
			styles.Value.Render(payload.Value),
		)
	}

	if payload.Value != "" {
		return styles.Value.Render(payload.Value)
	}

	lines := payload.Lines
	if len(lines) == 0 {
		raw := payload.Raw
		if raw == "" {
			raw = payload.Text
		}
		if raw != "" {
			lines = strings.Split(raw, "\n")
		}
	}

	if len(lines) == 0 {
		return ""
	}

	styled := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		styled = append(styled, styles.Value.Render(line))
	}

	return strings.Join(styled, "\n")
}

// LoadExternalPlugins registers executable scripts in the external directory.
func (pm *PluginManager) LoadExternalPlugins(dir string) error {
	if dir == "" {
		return nil
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.Mode()&0111 == 0 {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		pm.plugins[name] = &externalModule{
			name:    name,
			path:    path,
			timeout: pm.externalTimeout,
		}
	}

	return nil
}

var _ modules.Module = (*externalModule)(nil)
