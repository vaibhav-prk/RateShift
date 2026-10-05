// Package config provides loading and parsing
// application's YAML configuration.
package config

import (
	"os"

	"github.com/vaibhav-prk/RateShift/internal/limiter"
	"gopkg.in/yaml.v3"
)

type TenantConfig struct {
	TenantID  string                `yaml:"tenant_id"`
	APIKey    string                `yaml:"api_key"`
	RateLimit int                   `yaml:"rate_limit"`
	BurstSize int                   `yaml:"burst_size"`
	Algorithm limiter.AlgorithmType `yaml:"algorithm"`
}

type Config struct {
	Tenants map[string]TenantConfig `yaml:"tenants"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
