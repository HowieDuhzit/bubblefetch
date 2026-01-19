package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/collectors/local"
	"github.com/howieduhzit/bubblefetch/internal/collectors/remote"
	"github.com/howieduhzit/bubblefetch/internal/config"
	"github.com/howieduhzit/bubblefetch/internal/export"
	"github.com/howieduhzit/bubblefetch/internal/plugins"
	"github.com/howieduhzit/bubblefetch/internal/ui"
	"github.com/howieduhzit/bubblefetch/internal/ui/config_wizard"
	"github.com/howieduhzit/bubblefetch/internal/ui/modules"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
	"github.com/howieduhzit/bubblefetch/internal/whois"
)

var (
	configPath   = flag.String("config", "", "Path to config file (default: ~/.config/bubblefetch/config.yaml)")
	configPathS  = flag.String("c", "", "Alias for --config")
	themeName    = flag.String("theme", "", "Theme name to use")
	themeNameS   = flag.String("t", "", "Alias for --theme")
	remoteSys    = flag.String("remote", "", "Remote system IP/hostname to fetch info from (via SSH)")
	remoteSysS   = flag.String("r", "", "Alias for --remote")
	exportFmt    = flag.String("export", "", "Export format: json, yaml, or text")
	exportFmtS   = flag.String("e", "", "Alias for --export")
	pretty       = flag.Bool("pretty", true, "Pretty print JSON output (default: true)")
	prettyS      = flag.Bool("p", true, "Alias for --pretty")
	benchmark    = flag.Bool("benchmark", false, "Run benchmark mode")
	benchmarkS   = flag.Bool("b", false, "Alias for --benchmark")
	versionFlag  = flag.Bool("version", false, "Print version information")
	versionFlagS = flag.Bool("v", false, "Alias for --version")
	configWizard = flag.Bool("config-wizard", false, "Run interactive configuration wizard")
	configWizardS = flag.Bool("w", false, "Alias for --config-wizard")
	imageExport  = flag.String("image-export", "", "Export as image: png, svg, or html")
	imageOutput  = flag.String("image-output", "", "Image output path (default: bubblefetch.{format})")
	imageOutputS = flag.String("o", "", "Alias for --image-output")
	whoisTarget  = flag.String("who", "", "Domain scan (WHOIS + DNS records)")
	whoisTargetS = flag.String("W", "", "Alias for --who")
	whoisRaw     = flag.Bool("who-raw", false, "Include raw WHOIS output")
	whoisRawS    = flag.Bool("R", false, "Alias for --who-raw")
	helpFlag     = flag.Bool("help", false, "Show help message")
	helpFlagS    = flag.Bool("h", false, "Alias for --help")
)

const Version = "0.3.0"

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: bubblefetch [OPTIONS]

Options:
  -c, --config string         Path to config file (default: ~/.config/bubblefetch/config.yaml)
  -t, --theme string          Theme name to use (overrides config)
  -r, --remote string         Remote system IP/hostname to fetch info from (via SSH)
  -e, --export string         Export format: json, yaml, or text
  -p, --pretty                Pretty print JSON output (default: true)
  -b, --benchmark             Run benchmark mode (10 iterations)
  -w, --config-wizard         Run interactive configuration wizard
  --image-export string       Export as image: png, svg, or html
  -o, --image-output string   Image output path (default: bubblefetch.{format})
  -W, --who string            Domain scan (WHOIS + DNS records)
  -R, --who-raw               Include raw WHOIS output
  -v, --version               Print version information
  -h, --help                  Show help message

Notes:
  - If --image-export is omitted, the format is inferred from --image-output extension.
`)
	}
	flag.Parse()
	normalizeFlags()

	if *helpFlag {
		flag.Usage()
		return
	}

	if *versionFlag {
		fmt.Printf("bubblefetch v%s\n", Version)
		os.Exit(0)
	}

	// Run config wizard if requested
	if *configWizard {
		wizard := config_wizard.NewModel()
		p := tea.NewProgram(wizard)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running config wizard: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *whoisTarget != "" {
		runWhois(*whoisTarget, *whoisRaw)
		return
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override theme if specified
	if *themeName != "" {
		cfg.Theme = *themeName
	}

	// Override remote if specified
	if *remoteSys != "" {
		cfg.Remote = *remoteSys
	}

	// Load plugins
	pluginDir := cfg.PluginDir
	if pluginDir == "" {
		// Default plugin directory
		home, err := os.UserHomeDir()
		if err == nil {
			pluginDir = filepath.Join(home, ".config", "bubblefetch", "plugins")
		}
	}
	if pluginDir != "" {
		pm := plugins.NewPluginManager()
		if err := pm.LoadPlugins(pluginDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error loading plugins: %v\n", err)
		}
		modules.InitPlugins(pm)
	}

	// Handle image export mode
	if *imageExport != "" {
		runImageExport(cfg)
		return
	}

	// Handle export mode
	if *exportFmt != "" || *benchmark {
		runExportMode(cfg)
		return
	}

	// Create and run the Bubbletea program
	model := ui.NewModel(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func normalizeFlags() {
	if *configPathS != "" && *configPath == "" {
		*configPath = *configPathS
	}
	if *themeNameS != "" && *themeName == "" {
		*themeName = *themeNameS
	}
	if *remoteSysS != "" && *remoteSys == "" {
		*remoteSys = *remoteSysS
	}
	if *exportFmtS != "" && *exportFmt == "" {
		*exportFmt = *exportFmtS
	}
	if *imageOutputS != "" && *imageOutput == "" {
		*imageOutput = *imageOutputS
	}
	if *imageExport == "" && *imageOutput != "" {
		*imageExport = inferImageFormat(*imageOutput)
	}
	if *whoisTargetS != "" && *whoisTarget == "" {
		*whoisTarget = *whoisTargetS
	}
	if *whoisRawS && !*whoisRaw {
		*whoisRaw = true
	}
	if *benchmarkS && !*benchmark {
		*benchmark = true
	}
	if *configWizardS && !*configWizard {
		*configWizard = true
	}
	if *versionFlagS && !*versionFlag {
		*versionFlag = true
	}
	if *prettyS != *pretty {
		*pretty = *prettyS
	}
	if *helpFlagS && !*helpFlag {
		*helpFlag = true
	}
}

func inferImageFormat(outputPath string) string {
	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".png":
		return "png"
	case ".svg":
		return "svg"
	case ".html", ".htm":
		return "html"
	default:
		return ""
	}
}

func runExportMode(cfg *config.Config) {
	// Create collector
	var collector interface {
		Collect() (*collectors.SystemInfo, error)
	}

	if cfg.Remote != "" {
		collector = remote.New(cfg.Remote, cfg)
	} else {
		collector = local.New(cfg.EnablePublicIP)
	}

	// Collect system info
	var info *collectors.SystemInfo
	var err error

	if *benchmark {
		// Run benchmark
		const runs = 10
		var totalDuration time.Duration

		fmt.Printf("Running %d iterations...\n", runs)
		for i := 0; i < runs; i++ {
			start := time.Now()
			info, err = collector.Collect()
			duration := time.Since(start)
			totalDuration += duration
			fmt.Printf("Run %d: %v\n", i+1, duration)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error collecting info: %v\n", err)
				os.Exit(1)
			}
		}

		avgDuration := totalDuration / runs
		fmt.Printf("\nAverage: %v\n", avgDuration)
		fmt.Printf("Total: %v\n", totalDuration)
		return
	}

	info, err = collector.Collect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error collecting system info: %v\n", err)
		os.Exit(1)
	}

	// Export in requested format
	var output string
	switch *exportFmt {
	case "json":
		output, err = export.ToJSON(info, *pretty)
	case "yaml":
		output, err = export.ToYAML(info)
	case "text":
		output = export.ToText(info)
	default:
		fmt.Fprintf(os.Stderr, "Unknown export format: %s (use json, yaml, or text)\n", *exportFmt)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output)
}

func runImageExport(cfg *config.Config) {
	// Create collector
	var collector interface {
		Collect() (*collectors.SystemInfo, error)
	}

	if cfg.Remote != "" {
		collector = remote.New(cfg.Remote, cfg)
	} else {
		collector = local.New(cfg.EnablePublicIP)
	}

	// Collect system info
	info, err := collector.Collect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error collecting system info: %v\n", err)
		os.Exit(1)
	}

	// Create image exporter
	exporter, err := export.NewImageExporter(info, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating image exporter: %v\n", err)
		os.Exit(1)
	}

	// Determine output path
	outputPath := *imageOutput
	if outputPath == "" {
		outputPath = "bubblefetch." + *imageExport
	}

	// Export based on format
	switch *imageExport {
	case "png":
		err = exporter.ToPNG(outputPath)
	case "svg":
		err = exporter.ToSVG(outputPath)
	case "html":
		err = exporter.ToHTML(outputPath)
	default:
		fmt.Fprintf(os.Stderr, "Unknown image format: %s (use png, svg, or html)\n", *imageExport)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error exporting image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully exported to %s\n", outputPath)
}

func runWhois(target string, includeRaw bool) {
	cfg, _ := config.Load(*configPath)
	if cfg == nil {
		cfg = config.NewDefault()
	}
	if *themeName != "" {
		cfg.Theme = *themeName
	}

	thm, err := theme.Load(cfg.Theme)
	if err != nil {
		thm, _ = theme.Load("default")
	}

	result, err := whois.Lookup(target, includeRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(formatWhois(result, thm))
}

func formatWhois(result whois.Result, thm *theme.Theme) string {
	styles := thm.GetStyles()
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(thm.Colors.Accent))
	subhead := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(thm.Colors.Primary))
	separator := styles.Separator.Render(": ")
	icon := func(value string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(thm.Colors.Secondary)).Render(value)
	}

	var b strings.Builder
	b.WriteString(header.Render(icon("󰗇 ") + "Domain Scan"))
	b.WriteString("\n")
	b.WriteString(styles.Label.Render(icon("󰮱 ") + "Target"))
	b.WriteString(separator)
	b.WriteString(styles.Value.Render(result.Target))
	b.WriteString("\n\n")

	b.WriteString(subhead.Render(icon("󰌪 ") + "WHOIS"))
	b.WriteString("\n")
	if len(result.Whois.Fields) == 0 {
		b.WriteString(styles.Value.Render("  (no structured WHOIS fields found)"))
		b.WriteString("\n")
	} else {
		for _, field := range result.Whois.Fields {
			b.WriteString("  ")
			b.WriteString(styles.Label.Render(field.Label))
			b.WriteString(separator)
			b.WriteString(styles.Value.Render(field.Value))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(subhead.Render(icon("󰇋 ") + "DNS"))
	b.WriteString("\n")
	writeList(&b, styles, separator, "A", result.DNS.A)
	writeList(&b, styles, separator, "AAAA", result.DNS.AAAA)
	if result.DNS.CNAME != "" {
		writeLine(&b, styles, separator, "CNAME", result.DNS.CNAME)
	}
	writeList(&b, styles, separator, "MX", result.DNS.MX)
	writeList(&b, styles, separator, "NS", result.DNS.NS)
	writeList(&b, styles, separator, "PTR", result.DNS.PTR)
	if len(result.DNS.TXT) > 0 {
		b.WriteString("  ")
		b.WriteString(styles.Label.Render("TXT"))
		b.WriteString(":\n")
		for _, record := range result.DNS.TXT {
			b.WriteString("    - ")
			b.WriteString(styles.Value.Render(record))
			b.WriteString("\n")
		}
	}

	if result.Whois.Raw != "" {
		b.WriteString("\n")
		b.WriteString(subhead.Render(icon("󰡯 ") + "RAW WHOIS"))
		b.WriteString("\n")
		b.WriteString(result.Whois.Raw)
		b.WriteString("\n")
	}

	return styles.Border.Render(b.String())
}

func writeLine(b *strings.Builder, styles theme.Styles, separator, label, value string) {
	if value == "" {
		return
	}
	b.WriteString("  ")
	b.WriteString(styles.Label.Render(label))
	b.WriteString(separator)
	b.WriteString(styles.Value.Render(value))
	b.WriteString("\n")
}

func writeList(b *strings.Builder, styles theme.Styles, separator, label string, values []string) {
	if len(values) == 0 {
		return
	}
	writeLine(b, styles, separator, label, strings.Join(values, ", "))
}
