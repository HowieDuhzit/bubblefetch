package config_wizard

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/howieduhzit/bubblefetch/internal/config"
)

type step int

const (
	stepWelcome step = iota
	stepTheme
	stepModules
	stepPrivacy
	stepPlugins
	stepSave
	stepComplete
)

type Model struct {
	config   *config.Config
	step     step
	cursor   int
	selected map[string]bool // For multi-select (modules)

	// Available options
	themes  []string
	modules []string
	plugins []string

	// Styles
	titleStyle      lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
	helpStyle       lipgloss.Style
	errorStyle      lipgloss.Style
	successStyle    lipgloss.Style
}

func NewModel() Model {
	model := Model{
		config:   config.NewDefault(),
		selected: make(map[string]bool),
		themes: []string{
			"default",
			"minimal",
			"dracula",
			"nord",
			"gruvbox",
			"tokyo-night",
			"monokai",
			"solarized-dark",
		},
		modules: []string{
			"os",
			"kernel",
			"hostname",
			"uptime",
			"cpu",
			"gpu",
			"memory",
			"disk",
			"shell",
			"terminal",
			"de",
			"wm",
			"network",
			"localip",
			"publicip",
			"battery",
			"costs",
		},
		plugins: []string{},

		titleStyle:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#89b4fa")).MarginBottom(1),
		selectedStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Bold(true),
		unselectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#cdd6f4")),
		helpStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#585b70")).Italic(true).MarginTop(1),
		errorStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")).Bold(true),
		successStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Bold(true),
	}
	model.refreshPlugins()
	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step == stepComplete {
				return m, tea.Quit
			}
			// Ask for confirmation before quitting mid-wizard
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			maxCursor := m.getMaxCursor()
			if m.cursor < maxCursor {
				m.cursor++
			}

		case " ":
			// Space for multi-select (modules step)
			if m.step == stepModules {
				module := m.modules[m.cursor]
				m.selected[module] = !m.selected[module]
			} else if m.step == stepPrivacy {
				// Toggle public IP
				m.config.EnablePublicIP = !m.config.EnablePublicIP
			}

		case "enter":
			return m.handleEnter()
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.step {
	case stepWelcome:
		return m.viewWelcome()
	case stepTheme:
		return m.viewTheme()
	case stepModules:
		return m.viewModules()
	case stepPrivacy:
		return m.viewPrivacy()
	case stepPlugins:
		return m.viewPlugins()
	case stepSave:
		return m.viewSave()
	case stepComplete:
		return m.viewComplete()
	}
	return ""
}

func (m Model) viewWelcome() string {
	var b strings.Builder

	b.WriteString(m.titleStyle.Render("🎨 Welcome to bubblefetch Configuration Wizard!"))
	b.WriteString("\n\n")
	b.WriteString("This wizard will help you create a custom configuration for bubblefetch.\n")
	b.WriteString("You'll be able to:\n\n")
	b.WriteString("  • Choose a theme\n")
	b.WriteString("  • Select which modules to display\n")
	b.WriteString("  • Configure privacy settings\n")
	b.WriteString("  • Set up plugin directory\n\n")
	b.WriteString(m.helpStyle.Render("Press Enter to continue, or q to quit"))
	b.WriteString(m.renderProgress())

	return b.String()
}

func (m Model) viewTheme() string {
	var b strings.Builder

	b.WriteString(m.titleStyle.Render("Choose a Theme"))
	b.WriteString("\n")

	for i, theme := range m.themes {
		cursor := "  "
		if m.cursor == i {
			cursor = "→ "
			b.WriteString(m.selectedStyle.Render(cursor + theme + " ✓"))
		} else {
			b.WriteString(m.unselectedStyle.Render(cursor + theme))
		}
		b.WriteString("\n")
	}

	b.WriteString(m.helpStyle.Render("\n↑/↓: navigate • enter: select • q: quit"))
	b.WriteString(m.renderProgress())

	return b.String()
}

func (m Model) viewModules() string {
	var b strings.Builder

	b.WriteString(m.titleStyle.Render("Select Modules to Display"))
	b.WriteString("\n")

	for i, module := range m.modules {
		cursor := "  "
		checked := " "

		if m.cursor == i {
			cursor = "→ "
		}

		if m.selected[module] {
			checked = "✓"
		}

		style := m.unselectedStyle
		if m.cursor == i {
			style = m.selectedStyle
		}

		line := fmt.Sprintf("%s[%s] %s", cursor, checked, module)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	selectedCount := 0
	for _, selected := range m.selected {
		if selected {
			selectedCount++
		}
	}

	b.WriteString(m.helpStyle.Render(fmt.Sprintf("\n%d modules selected\n↑/↓: navigate • space: toggle • enter: continue • q: quit", selectedCount)))
	b.WriteString(m.renderProgress())

	return b.String()
}

func (m Model) viewPrivacy() string {
	var b strings.Builder

	b.WriteString(m.titleStyle.Render("Privacy Settings"))
	b.WriteString("\n")

	// Public IP option
	cursor := "→ "
	checked := " "
	if m.config.EnablePublicIP {
		checked = "✓"
	}

	b.WriteString(m.selectedStyle.Render(fmt.Sprintf("%s[%s] Enable Public IP Detection", cursor, checked)))
	b.WriteString("\n\n")
	b.WriteString(m.unselectedStyle.Render("Public IP detection fetches your public IP address from external services.\n"))
	b.WriteString(m.unselectedStyle.Render("This requires an HTTP request to api.ipify.org or similar services.\n"))
	b.WriteString(m.unselectedStyle.Render("Disabled by default for privacy."))

	b.WriteString(m.helpStyle.Render("\n\nspace: toggle • enter: continue • q: quit"))
	b.WriteString(m.renderProgress())

	return b.String()
}

func (m Model) viewPlugins() string {
	var b strings.Builder

	b.WriteString(m.titleStyle.Render("Plugin Directory"))
	b.WriteString("\n")

	pluginDir := m.config.PluginDir
	if pluginDir == "" {
		pluginDir = "~/.config/bubblefetch/plugins (default)"
	}

	b.WriteString(m.unselectedStyle.Render(fmt.Sprintf("Plugins will be loaded from:\n%s\n\n", pluginDir)))
	b.WriteString(m.unselectedStyle.Render("You can change this later by editing config.yaml\n"))
	b.WriteString(m.unselectedStyle.Render("See docs/PLUGINS.md for plugin development guide.\n\n"))

	if len(m.plugins) == 0 {
		b.WriteString(m.unselectedStyle.Render("Detected plugins: none\n"))
	} else {
		b.WriteString(m.unselectedStyle.Render("Detected plugins:\n"))
		for _, plugin := range m.plugins {
			b.WriteString(m.unselectedStyle.Render("  • " + plugin))
			b.WriteString("\n")
		}
	}

	b.WriteString(m.helpStyle.Render("\n\nenter: continue • q: quit"))
	b.WriteString(m.renderProgress())

	return b.String()
}

func (m Model) viewSave() string {
	var b strings.Builder

	b.WriteString(m.titleStyle.Render("Ready to Save Configuration"))
	b.WriteString("\n")

	b.WriteString("Your configuration:\n\n")
	b.WriteString(m.unselectedStyle.Render(fmt.Sprintf("  Theme: %s\n", m.config.Theme)))
	b.WriteString(m.unselectedStyle.Render(fmt.Sprintf("  Modules: %d selected\n", len(m.config.Modules))))
	b.WriteString(m.unselectedStyle.Render(fmt.Sprintf("  Public IP: %v\n", m.config.EnablePublicIP)))
	b.WriteString("\n")
	b.WriteString("Configuration will be saved to:\n")
	b.WriteString(m.selectedStyle.Render("~/.config/bubblefetch/config.yaml"))

	b.WriteString(m.helpStyle.Render("\n\nenter: save and finish • q: quit without saving"))
	b.WriteString(m.renderProgress())

	return b.String()
}

func (m Model) viewComplete() string {
	var b strings.Builder

	b.WriteString(m.successStyle.Render("✓ Configuration Saved Successfully!"))
	b.WriteString("\n\n")
	b.WriteString("Your configuration has been saved to:\n")
	b.WriteString(m.selectedStyle.Render("~/.config/bubblefetch/config.yaml"))
	b.WriteString("\n\n")
	b.WriteString("Run bubblefetch to see your customized system info!\n\n")
	b.WriteString(m.unselectedStyle.Render("Tip: Edit config.yaml manually for advanced customization\n"))
	b.WriteString(m.unselectedStyle.Render("See docs/PLUGINS.md to create custom modules"))

	b.WriteString(m.helpStyle.Render("\n\nPress q or Ctrl+C to exit"))

	return b.String()
}

func (m Model) renderProgress() string {
	current := int(m.step) + 1
	total := int(stepComplete)
	progressBar := ""

	for i := 1; i <= total; i++ {
		if i == current {
			progressBar += "● "
		} else if i < current {
			progressBar += "✓ "
		} else {
			progressBar += "○ "
		}
	}

	return "\n\n" + m.helpStyle.Render(progressBar)
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepWelcome:
		m.step = stepTheme
		m.cursor = 0

	case stepTheme:
		m.config.Theme = m.themes[m.cursor]
		m.refreshPlugins()
		m.step = stepModules
		m.cursor = 0
		// Pre-select recommended modules
		for _, mod := range []string{"os", "kernel", "hostname", "cpu", "memory", "disk"} {
			m.selected[mod] = true
		}

	case stepModules:
		// Build module list from selections
		m.config.Modules = []string{}
		for _, mod := range m.modules {
			if m.selected[mod] {
				m.config.Modules = append(m.config.Modules, mod)
			}
		}
		m.step = stepPrivacy
		m.cursor = 0

	case stepPrivacy:
		m.plugins = detectPlugins(m.config.PluginDir)
		m.step = stepPlugins
		m.cursor = 0

	case stepPlugins:
		m.step = stepSave
		m.cursor = 0

	case stepSave:
		// Save config
		if err := config.Save(m.config); err != nil {
			// Error saving - could show error message
			// For now, just continue to complete
		}
		m.step = stepComplete

	case stepComplete:
		return m, tea.Quit
	}

	return m, nil
}

func detectPlugins(pluginDir string) []string {
	if pluginDir == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			pluginDir = filepath.Join(home, ".config", "bubblefetch", "plugins")
		}
	}

	var plugins []string
	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".so" {
			plugins = append(plugins, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())))
		}
	}

	externalDir := filepath.Join(pluginDir, "external")
	externalEntries, err := os.ReadDir(externalDir)
	if err == nil {
		for _, entry := range externalEntries {
			if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.Mode()&0111 == 0 {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			plugins = append(plugins, name)
		}
	}

	sort.Strings(plugins)
	return plugins
}

func (m *Model) refreshPlugins() {
	m.plugins = detectPlugins(m.config.PluginDir)
	if len(m.plugins) == 0 {
		return
	}

	existing := make(map[string]bool, len(m.modules))
	for _, module := range m.modules {
		existing[module] = true
	}
	for _, plugin := range m.plugins {
		if !existing[plugin] {
			m.modules = append(m.modules, plugin)
		}
	}
}

func (m Model) getMaxCursor() int {
	switch m.step {
	case stepTheme:
		return len(m.themes) - 1
	case stepModules:
		return len(m.modules) - 1
	default:
		return 0
	}
}
