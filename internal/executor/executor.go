package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/glowfi/autost/internal/config"
	"github.com/google/uuid"
)

// Executor runs commands and scripts
type Executor struct{}

// New creates a new executor
func New() *Executor {
	return &Executor{}
}

func (e *Executor) Run(ctx context.Context, cfg config.Config) error {
	if err := e.setEnvVars(cfg.Env); err != nil {
		return err
	}

	for _, cmd := range cfg.Startup {
		if err := e.executeCommand(ctx, cmd, cfg.Interpreter); err != nil {
			return err
		}
	}

	for _, script := range cfg.Scripts {
		if err := e.executeScript(ctx, cfg.Interpreter, script, uuid.New().String()); err != nil {
			return err
		}
	}

	return nil
}

// executeCommand runs a simple command
func (e *Executor) executeCommand(ctx context.Context, command string, interpreter string) error {
	cmd := exec.CommandContext(ctx, interpreter, "-c", command)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	return nil
}

// executeScript runs an inline script with interpreter
func (e *Executor) executeScript(ctx context.Context, interpreter string, script string, name string) error {
	// Create temp file
	tmpFile, err := os.CreateTemp("/tmp", fmt.Sprintf("hotkeysd-*%s", name))
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	tmpPath := tmpFile.Name()

	// Write script content
	content := script
	if _, err := tmpFile.WriteString(content); err != nil {
		if err := tmpFile.Close(); err != nil {
			return err
		}
		if err := os.Remove(tmpPath); err != nil {
			return err
		}
		return fmt.Errorf("write script: %w", err)
	}
	tmpFile.Close()

	// Make executable
	os.Chmod(tmpPath, 0o700)

	// Execute
	cmd := exec.CommandContext(ctx, interpreter, tmpPath)

	if err := cmd.Start(); err != nil {
		if err := os.Remove(tmpPath); err != nil {
			return err
		}
		return fmt.Errorf("start script: %w", err)
	}

	return nil
}

// sets enviroment variable
func (e *Executor) setEnvVars(envVars []config.EnvVar) error {
	for _, env := range envVars {
		if err := os.Setenv(env.Key, env.Value); err != nil {
			return fmt.Errorf("setting %s: %w", env.Key, err)
		}
	}
	return nil
}
