package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/glowfi/autost/internal/config"
	"github.com/glowfi/autost/internal/executor"
)

func main() {
	// Find config file
	configPath, err := findConfig()
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		fmt.Printf("\nReceived %s, shutting down...\n", sig)
		cancel()
	}()

	fmt.Printf("autost: loading config from %s\n", configPath)

	ex := executor.NewExecutor()
	if err := ex.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("autost: startup complete")

	// Wait for shutdown signal
	<-ctx.Done()
}

func findConfig() (string, error) {
	// Check XDG_CONFIG_HOME
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}

	paths := []string{
		filepath.Join(configHome, "autost", "config.yaml"),
		filepath.Join(configHome, "autost", "config.yml"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("config not found")
}
