package theme

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name      string `json:"name"`
	Colors    Colors `json:"colors"`
	ASCII     string `json:"ascii"`
	Layout    Layout `json:"layout"`
	asciiAuto bool
}

type Colors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Accent     string `json:"accent"`
	Label      string `json:"label"`
	Value      string `json:"value"`
	Border     string `json:"border"`
	Background string `json:"background"`
}

type Layout struct {
	ShowASCII   bool   `json:"show_ascii"`
	ASCIIWidth  int    `json:"ascii_width"`
	Separator   string `json:"separator"`
	Padding     int    `json:"padding"`
	BorderStyle string `json:"border_style"`
}

// Styles contains lipgloss styles for the theme
type Styles struct {
	Label     lipgloss.Style
	Value     lipgloss.Style
	Separator lipgloss.Style
	ASCII     lipgloss.Style
	Border    lipgloss.Style
	Container lipgloss.Style
}

// Load loads a theme by name from the themes directory
func Load(name string) (*Theme, error) {
	// Try user config directory first
	home, err := os.UserHomeDir()
	if err == nil {
		userThemePath := filepath.Join(home, ".config", "bubblefetch", "themes", name+".json")
		if theme, err := loadThemeFile(userThemePath); err == nil {
			return theme, nil
		}
	}

	// Try local themes directory
	themePath := filepath.Join("themes", name+".json")
	if theme, err := loadThemeFile(themePath); err == nil {
		return theme, nil
	}

	// Fall back to default theme
	return defaultTheme(), nil
}

func loadThemeFile(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, err
	}

	// Auto-detect ASCII art if not set or set to "auto"
	theme.asciiAuto = theme.ASCII == "" || theme.ASCII == "auto"
	if theme.asciiAuto {
		detectedOS := DetectOS()
		theme.ASCII = GetASCIIArt(detectedOS)
	}

	return &theme, nil
}

func defaultTheme() *Theme {
	detectedOS := DetectOS()
	return &Theme{
		Name: "default",
		Colors: Colors{
			Primary:    "#89b4fa",
			Secondary:  "#cba6f7",
			Accent:     "#f38ba8",
			Label:      "#f9e2af",
			Value:      "#a6e3a1",
			Border:     "#585b70",
			Background: "#1e1e2e",
		},
		ASCII:     GetASCIIArt(detectedOS),
		asciiAuto: true,
		Layout: Layout{
			ShowASCII:   true,
			ASCIIWidth:  30,
			Separator:   " ",
			Padding:     2,
			BorderStyle: "rounded",
		},
	}
}

// ApplyAutoASCII updates ASCII art when the theme is set to auto.
func (t *Theme) ApplyAutoASCII(osName string) {
	if t == nil || !t.asciiAuto || osName == "" {
		return
	}
	t.ASCII = GetASCIIArt(osName)
}

// GetStyles creates lipgloss styles from the theme
func (t *Theme) GetStyles() Styles {
	var borderStyle lipgloss.Border
	switch t.Layout.BorderStyle {
	case "rounded":
		borderStyle = lipgloss.RoundedBorder()
	case "double":
		borderStyle = lipgloss.DoubleBorder()
	case "thick":
		borderStyle = lipgloss.ThickBorder()
	default:
		borderStyle = lipgloss.NormalBorder()
	}

	return Styles{
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Label)).
			Bold(true),
		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Value)),
		Separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Primary)),
		ASCII: lipgloss.NewStyle().
			Foreground(lipgloss.Color(t.Colors.Primary)),
		Border: lipgloss.NewStyle().
			BorderStyle(borderStyle).
			BorderForeground(lipgloss.Color(t.Colors.Border)).
			Padding(0, t.Layout.Padding),
		Container: lipgloss.NewStyle().
			Padding(1, 2),
	}
}
