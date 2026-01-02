package app

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type HabitConfig struct {
	ID      string `yaml:"id"`
	Label   string `yaml:"label"`
	Column  string `yaml:"column"`
	Default bool   `yaml:"default"`
	Help    string `yaml:"help"`
}

type AuthConfig struct {
	Mode    string `yaml:"mode"`
	KeyPath string `yaml:"key_path"`
}

type Config struct {
	SpreadsheetID string        `yaml:"spreadsheet_id"`
	SheetName     string        `yaml:"sheet_name"`
	Timezone      string        `yaml:"timezone"`
	Auth          AuthConfig    `yaml:"auth"`
	Habits        []HabitConfig `yaml:"habits"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(home, ".config", "gohabits", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Timezone == "" {
		cfg.Timezone = "Local"
	}

	return &cfg, nil
}
