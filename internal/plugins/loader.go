package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	"github.com/howieduhzit/bubblefetch/internal/collectors"
	"github.com/howieduhzit/bubblefetch/internal/ui/modules"
	"github.com/howieduhzit/bubblefetch/internal/ui/theme"
)

// PluginModule wraps a plugin's render function to implement the Module interface
type PluginModule struct {
	name   string
	render func(*collectors.SystemInfo, theme.Styles) string
}

func (p *PluginModule) Name() string {
	return p.name
}

func (p *PluginModule) Render(info *collectors.SystemInfo, styles theme.Styles) string {
	return p.render(info, styles)
}

// PluginManager manages loading and accessing plugins
type PluginManager struct {
	plugins map[string]modules.Module
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]modules.Module),
	}
}

// LoadPlugins loads all .so files from the specified directory
func (pm *PluginManager) LoadPlugins(pluginDir string) error {
	// Check if directory exists
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		// Directory doesn't exist, silently return (no plugins to load)
		return nil
	}

	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".so" {
			continue
		}

		pluginPath := filepath.Join(pluginDir, entry.Name())
		if err := pm.LoadPlugin(pluginPath); err != nil {
			// Log warning but continue loading other plugins
			fmt.Fprintf(os.Stderr, "Warning: failed to load plugin %s: %v\n", entry.Name(), err)
			continue
		}
	}

	return nil
}

// LoadPlugin loads a single plugin from the specified path
func (pm *PluginManager) LoadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look for required symbols
	nameSymbol, err := p.Lookup("ModuleName")
	if err != nil {
		return fmt.Errorf("plugin missing ModuleName: %w", err)
	}

	renderSymbol, err := p.Lookup("Render")
	if err != nil {
		return fmt.Errorf("plugin missing Render function: %w", err)
	}

	// Type assert to expected types
	name, ok := nameSymbol.(*string)
	if !ok {
		return fmt.Errorf("ModuleName is not *string")
	}

	renderFunc, ok := renderSymbol.(func(*collectors.SystemInfo, theme.Styles) string)
	if !ok {
		return fmt.Errorf("Render is not the correct function type")
	}

	// Register plugin
	pm.plugins[*name] = &PluginModule{
		name:   *name,
		render: renderFunc,
	}

	return nil
}

// GetPlugin retrieves a plugin by name
func (pm *PluginManager) GetPlugin(name string) (modules.Module, bool) {
	mod, ok := pm.plugins[name]
	return mod, ok
}

// ListPlugins returns a list of all loaded plugin names
func (pm *PluginManager) ListPlugins() []string {
	names := make([]string, 0, len(pm.plugins))
	for name := range pm.plugins {
		names = append(names, name)
	}
	return names
}
