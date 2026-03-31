package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
}

// Load loads the configuration from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	// Override with env if present
	if envDSN := os.Getenv("DB_DSN"); envDSN != "" {
		cfg.DB.DSN = envDSN
	}

	return &cfg, nil
}
