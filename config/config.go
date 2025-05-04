// SPDX-FileCopyrightText: 2025 Sidings Media <contact@sidingsmedia.com>
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

type Target struct {
    Host string `yaml:"host"`
    Interface string `yaml:"interface"`
    TTL int8 `yaml:"ttl"`
    Size int `yaml:"size"`
    Count int `yaml:"count"`
}

type Config struct {
    Targets []Target `yaml:"targets"`
    DefaultTTL int8 `yaml:"default_ttl"`
    DefaultSize int `yaml:"default_size"`
    DefaultCount int `yaml:"default_count"`
}

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

    logger.Debug("Loaded configuration from file", "configuration", config)
    if len(config.Targets) < 1 {
        logger.Warn("No targets specified. Did you forget to set targets in your config file?")
    }

    return &config, nil
}
