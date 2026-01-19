package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/collectors/local"
	"github.com/howieduhzit/bubblefetch/internal/collectors/remote"
	"github.com/howieduhzit/bubblefetch/internal/config"
	"github.com/howieduhzit/bubblefetch/internal/ui/modules"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

type Model struct {
	config    *config.Config
	theme     *theme.Theme
	styles    theme.Styles
	sysInfo   *collectors.SystemInfo
	collector collectors.Collector
	err       error
	ready     bool
}

type collectMsg struct {
	info *collectors.SystemInfo
	err  error
}

func NewModel(cfg *config.Config) Model {
	// Load theme
	thm, err := theme.Load(cfg.Theme)
	if err != nil {
		thm, _ = theme.Load("default")
	}

	// Create collector based on config
	var collector collectors.Collector
	if cfg.Remote != "" {
		collector = remote.New(cfg.Remote, cfg)
	} else {
		collector = local.New(cfg.EnablePublicIP)
	}

	return Model{
		config:    cfg,
		theme:     thm,
		styles:    thm.GetStyles(),
		collector: collector,
		ready:     false,
	}
}

func (m Model) Init() tea.Cmd {
	return collectSystemInfo(m.collector)
}

func collectSystemInfo(c collectors.Collector) tea.Cmd {
	return func() tea.Msg {
		info, err := c.Collect()
		return collectMsg{info: info, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c", "esc"))):
			return m, tea.Quit
		}

	case collectMsg:
		m.sysInfo = msg.info
		m.err = msg.err
		m.ready = true
	}

	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "Collecting system information..."
	}

	if m.err != nil {
		return m.styles.Value.Render("Error: " + m.err.Error())
	}

	result := Render(m.config, m.sysInfo, nil)

	// Add help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.theme.Colors.Border)).Faint(true)
	help := helpStyle.Render("\nPress q to quit")

	return result + help
}
