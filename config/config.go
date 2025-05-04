// SPDX-FileCopyrightText: 2025 Sidings Media <contact@sidingsmedia.com>
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
)

// Taken from https://github.com/prometheus-community/pro-bing/blob/85df87ee97d5a448f5bc5c2ccc6f43d54e68b0cd/ping.go#L77C2-L78C37
const minSize = 8 + len(uuid.UUID{})

type Target struct {
	Host      string `yaml:"host"`
	Interface string `yaml:"interface"`
	TTL       int    `yaml:"ttl"`
	Size      int    `yaml:"size"`
	Count     int    `yaml:"count"`
	Interval  int    `yaml:"interval"`
}

type Config struct {
	Targets         []Target `yaml:"targets"`
	DefaultTTL      int      `yaml:"default_ttl"`
	DefaultSize     int      `yaml:"default_size"`
	DefaultCount    int      `yaml:"default_count"`
	DefaultInterval int      `yaml:"default_interval"`
	Timeout         int      `yaml:"timeout"`
}

// Parse and validate configuration.
func ParseConfig(file string, logger *slog.Logger) (*Config, error) {
	logger.Debug(fmt.Sprintf("Loading configuration from %s", file))

	dat, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := Config{}

	if err := yaml.Unmarshal(dat, &config); err != nil {
		return nil, err
	}

	setDefaults(&config)
	if err := validateSize(&config); err != nil {
		return nil, err
	}

	logger.Debug("Loaded configuration from file", "configuration", config)
	if len(config.Targets) < 1 {
		logger.Warn("No targets specified. Did you forget to set targets in your config file?")
	}

	return &config, nil
}

// Set default values where not already populated
func setDefaults(config *Config) {
	if config.DefaultCount == 0 {
		config.DefaultCount = 1
	}

	if config.DefaultSize == 0 {
		config.DefaultSize = minSize
	}

	if config.DefaultInterval == 0 {
		config.DefaultInterval = 1000
	}

	if config.Timeout == 0 {
		config.Timeout = 5000 // 5 Seconds
	}

	for i := range config.Targets {
		target := &config.Targets[i]
		if target.Count == 0 {
			target.Count = config.DefaultCount
		}

		if target.Size == 0 {
			target.Size = config.DefaultSize
		}

		if target.TTL == 0 {
			target.TTL = config.DefaultTTL
		}

		if target.Interval == 0 {
			target.Interval = config.DefaultInterval
		}
	}
}

// Ensure that the packet size is not invalid
func validateSize(config *Config) error {
	if config.DefaultSize < minSize {
		return fmt.Errorf("default packet size %d is less than minimum %d", config.DefaultSize, minSize)
	}

	for _, target := range config.Targets {
		if target.Size < minSize {
			return fmt.Errorf("packet size %d for target %s is less than minimum %d", target.Size, target.Host, minSize)
		}
	}
	return nil
}
