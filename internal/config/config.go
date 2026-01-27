package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Env         map[string]string `yaml:"env"`
	Interpreter string            `yaml:"interpreter"`
	Startup     []string          `yaml:"startup"`
	Scripts     []string          `yaml:"scripts"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
