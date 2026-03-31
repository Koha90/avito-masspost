// Package config loads application configuration from TOML files.
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config contains full application configuration.
type Config struct {
	Avito Avito `toml:"avito"`
	Feed  Feed  `toml:"feed"`
}

// Avito contains Avito API settings.
type Avito struct {
	BaseURL      string `toml:"base_url"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	UserID       string `toml:"user_id"`
}

// Feed contains feed output settings.
type Feed struct {
	OutputPath string `toml:"output_path"`
}

// Path returns a config file path from CLI flag, environment or default value.
func Path() string {
	var path string

	flag.StringVar(&path, "config", "", "path to TOML config file")
	flag.Parse()

	if path != "" {
		return filepath.Clean(path)
	}

	if env := os.Getenv("AVITO_CONFIG"); env != "" {
		return filepath.Clean(env)
	}

	return "config.toml"
}

// Load reads configuration from a TOML file.
func Load(path string) (Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("cannot read config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks required config fields.
func (c Config) Validate() error {
	if c.Avito.BaseURL == "" {
		return errors.New("avito.base_url is required")
	}
	if c.Avito.ClientID == "" {
		return errors.New("avito.client_id is required")
	}
	if c.Avito.ClientSecret == "" {
		return errors.New("avito.client_secret is required")
	}
	if c.Avito.UserID == "" {
		return errors.New("avito.user_id is required")
	}
	if c.Feed.OutputPath == "" {
		return errors.New("feed.output_path is required")
	}

	return nil
}
