package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/collectors/local"
	"github.com/howieduhzit/bubblefetch/internal/collectors/remote"
	"github.com/howieduhzit/bubblefetch/internal/config"
	"github.com/howieduhzit/bubblefetch/internal/export"
	"github.com/howieduhzit/bubblefetch/internal/plugins"
	"github.com/howieduhzit/bubblefetch/internal/ui"
	"github.com/howieduhzit/bubblefetch/internal/ui/config_wizard"
	"github.com/howieduhzit/bubblefetch/internal/ui/modules"
)

var (
	configPath   = flag.String("config", "", "Path to config file (default: ~/.config/bubblefetch/config.yaml)")
	themeName    = flag.String("theme", "", "Theme name to use")
	remoteSys    = flag.String("remote", "", "Remote system IP/hostname to fetch info from (via SSH)")
	exportFmt    = flag.String("export", "", "Export format: json, yaml, or text")
	pretty       = flag.Bool("pretty", true, "Pretty print JSON output (default: true)")
	benchmark    = flag.Bool("benchmark", false, "Run benchmark mode")
	versionFlag  = flag.Bool("version", false, "Print version information")
	configWizard = flag.Bool("config-wizard", false, "Run interactive configuration wizard")
	imageExport  = flag.String("image-export", "", "Export as image: png, svg, or html")
	imageOutput  = flag.String("image-output", "", "Image output path (default: bubblefetch.{format})")
)

const Version = "0.3.0"

func main() {
	flag.Parse()

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
