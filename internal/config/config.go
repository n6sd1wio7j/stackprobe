package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Endpoint represents a single service endpoint to health-check.
type Endpoint struct {
	Name    string        `json:"name"`
	URL     string        `json:"url"`
	Timeout time.Duration `json:"timeout"`
}

// Config holds the full stackprobe configuration.
type Config struct {
	ListenAddr string        `json:"listen_addr"`
	Interval   time.Duration `json:"interval"`
	Endpoints  []Endpoint    `json:"endpoints"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		ListenAddr: ":8080",
		Interval:   30 * time.Second,
		Endpoints:  []Endpoint{},
	}
}

// Load reads a JSON config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required fields are present and values are sane.
func (c *Config) validate() error {
	if c.ListenAddr == "" {
		return fmt.Errorf("config: listen_addr must not be empty")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("config: interval must be positive")
	}
	for i, ep := range c.Endpoints {
		if ep.Name == "" {
			return fmt.Errorf("config: endpoint[%d]: name must not be empty", i)
		}
		if ep.URL == "" {
			return fmt.Errorf("config: endpoint[%d] %q: url must not be empty", i, ep.Name)
		}
		if ep.Timeout <= 0 {
			c.Endpoints[i].Timeout = 5 * time.Second
		}
	}
	return nil
}
