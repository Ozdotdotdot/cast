package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	DefaultDevice      string            `toml:"default_device"`
	BlockUntilComplete bool              `toml:"block_until_complete"`
	Devices            map[string]Device `toml:"devices"`
	Groups             map[string]Group  `toml:"groups"`
}

type Device struct {
	IP           string   `toml:"ip"`
	Port         int      `toml:"port"`
	OriginalName string   `toml:"original_name"`
	Aliases      []string `toml:"aliases,omitempty"`
}

type Group struct {
	Devices []string `toml:"devices"`
}

func New() *Config {
	return &Config{
		BlockUntilComplete: true,
		Devices:            make(map[string]Device),
		Groups:             make(map[string]Group),
	}
}

func Path() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = os.Getenv("HOME")
	}
	return filepath.Join(configDir, "cast", "config.toml")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Initialize maps if nil
	if cfg.Devices == nil {
		cfg.Devices = make(map[string]Device)
	}
	if cfg.Groups == nil {
		cfg.Groups = make(map[string]Group)
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path := Path()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := toml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetDevice finds a device by name or alias
func (c *Config) GetDevice(nameOrAlias string) (Device, error) {
	// Check direct name match
	if d, ok := c.Devices[nameOrAlias]; ok {
		return d, nil
	}

	// Check aliases
	for _, d := range c.Devices {
		for _, alias := range d.Aliases {
			if alias == nameOrAlias {
				return d, nil
			}
		}
	}

	return Device{}, fmt.Errorf("device not found: %s", nameOrAlias)
}
