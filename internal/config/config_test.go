package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		cfgPath    string
		wantConfig Config
		wantErr    error
	}{
		{
			name:       "should return error when loading config with missing values from env variables :NEG",
			cfgPath:    "testdata/invalid_env_var.yaml",
			wantConfig: Config{},
			wantErr:    ErrEmptyValue,
		},
		{
			name:       "should return no error when loading empty config :POS",
			cfgPath:    "testdata/empty.yaml",
			wantConfig: Config{},
			wantErr:    nil,
		},
		{
			name:    "should load config :POS",
			cfgPath: "testdata/config.yaml",
			wantConfig: Config{
				Env: []EnvVar{
					{Key: "EDITOR", Value: "nvim"},
					{Key: "PATH", Value: "${PATH}:${HOME}/.local/bin"},
					{Key: "XDG_CONFIG_HOME", Value: "${HOME}/.config"},
				},
				Interpreter: "/bin/sh",
				Startup: []string{
					"dunst",
					"picom --backend glx",
				},
				Scripts: []string{
					"hour=$(date +%H)\necho \"${hour}\"\n",
					"echo \"hello\"\n",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, gotErr := LoadConfig(tt.cfgPath)

			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			assert.Equal(t, tt.wantConfig, gotConfig, "expect config to match")
		})
	}
}
