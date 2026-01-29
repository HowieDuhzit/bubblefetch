package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme                   string    `yaml:"theme"`
	Remote                  string    `yaml:"remote"`
	Modules                 []string  `yaml:"modules"`
	SSH                     SSHConfig `yaml:"ssh"`
	EnablePublicIP          bool      `yaml:"enable_public_ip"`
	PluginDir               string    `yaml:"plugin_dir"`
	ExternalModuleTimeoutMS int       `yaml:"external_module_timeout_ms"`
	BagsAPIKey              string    `yaml:"bags_api_key"` // Optional Bags.fm API key for enhanced Solana token data
}

type SSHConfig struct {
	User           string `yaml:"user"`
	Port           int    `yaml:"port"`
	KeyPath        string `yaml:"key_path"`
	KnownHostsPath string `yaml:"known_hosts_path"`
	SafeMode       bool   `yaml:"safe_mode"`
}

// Load loads configuration from the specified path or default location
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return defaultConfig(), nil
		}
		configPath = filepath.Join(home, ".config", "bubblefetch", "config.yaml")
	}

	// If config file doesn't exist, use defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return defaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Fill in defaults for missing fields
	if cfg.Theme == "" {
		cfg.Theme = "default"
	}
	if cfg.SSH.Port == 0 {
		cfg.SSH.Port = 22
	}
	if cfg.ExternalModuleTimeoutMS == 0 {
		cfg.ExternalModuleTimeoutMS = 250
	}
	if len(cfg.Modules) == 0 {
		cfg.Modules = defaultModules()
	}

	return &cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Theme:   "default",
		Modules: defaultModules(),
		SSH: SSHConfig{
			Port: 22,
		},
		ExternalModuleTimeoutMS: 250,
	}
}

// NewDefault creates a new default configuration
func NewDefault() *Config {
	return defaultConfig()
}

// Save saves the configuration to the default config file
func Save(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(home, ".config", "bubblefetch", "config.yaml")

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func defaultModules() []string {
	return []string{
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
		"localip",
		"battery",
	}
}
