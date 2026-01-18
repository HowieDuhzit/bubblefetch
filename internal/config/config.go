package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme   string   `yaml:"theme"`
	Remote  string   `yaml:"remote"`
	Modules []string `yaml:"modules"`
	SSH     SSHConfig `yaml:"ssh"`
}

type SSHConfig struct {
	User           string `yaml:"user"`
	Port           int    `yaml:"port"`
	KeyPath        string `yaml:"key_path"`
	KnownHostsPath string `yaml:"known_hosts_path"`
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
	}
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
