// Package config loads application configuration from TOML files.
package config

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config contains full application configuration.
type Config struct {
	Database  DatabaseConfig  `toml:"database"`
	Migration MigrationConfig `toml:"migration"`
	Avito     AvitoConfig     `toml:"avito"`
	Feed      FeedConfig      `toml:"feed"`
}

// DatabaseConfig contains PostgreSQL connection settings.
type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
	SSLMode  string `toml:"sslmode"`
}

// DSN returns a PostgreSQL connection string.
func (c DatabaseConfig) DSN() string {
	q := url.Values{}
	q.Set("sslmode", c.SSLMode)

	u := url.URL{
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:     c.Name,
		RawQuery: q.Encode(),
	}

	if c.User != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	return u.String()
}

// MigrationConfig contains migration settings.
type MigrationConfig struct {
	Path string `toml:"path"`
}

// AvitoConfig contains Avito API settings.
type AvitoConfig struct {
	BaseURL      string `toml:"base_url"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	UserID       string `toml:"user_id"`
}

// FeedConfig contains feed output settings.
type FeedConfig struct {
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

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks required config fields.
func (c Config) Validate() error {
	if c.Database.Host == "" {
		return errors.New("database.host is required")
	}
	if c.Database.Port <= 0 {
		return errors.New("database.port must be positive")
	}
	if c.Database.User == "" {
		return errors.New("database.user is required")
	}
	if c.Database.Name == "" {
		return errors.New("database.name is required")
	}
	if c.Database.SSLMode == "" {
		return errors.New("database.sslmode is required")
	}
	if c.Migration.Path == "" {
		return errors.New("migration.path is required")
	}
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
