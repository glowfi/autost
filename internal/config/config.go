package config

import (
	"errors"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

var ErrEmptyValue = errors.New("values of environment variables cannot be empty")

type EnvVar struct {
	Key   string
	Value string
}

func (e *EnvVar) UnmarshalYAML(data []byte) error {
	// Single key-value pair comes as a map
	var m map[string]string
	if err := yaml.Unmarshal(data, &m); err != nil {
		return err
	}

	var key string
	var val string

	for k, v := range m {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			return ErrEmptyValue
		}
		key = k
		val = v
		break
	}

	*e = EnvVar{
		Key:   key,
		Value: val,
	}

	return nil
}

func (e EnvVar) MarshalYAML() ([]byte, error) {
	m := map[string]string{e.Key: e.Value}
	return yaml.Marshal(m)
}

type Config struct {
	Env         []EnvVar `yaml:"env"`
	Interpreter string   `yaml:"interpreter"`
	Startup     []string `yaml:"startup"`
	Scripts     []string `yaml:"scripts"`
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
