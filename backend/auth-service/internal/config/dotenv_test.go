package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsDotEnvAndPreservesExportedEnvironment(t *testing.T) {
	tempDir := t.TempDir()
	dotEnv := "DOTENV_TEST_FROM_FILE=loaded\nDOTENV_TEST_PRECEDENCE=from-file\n"
	if err := os.WriteFile(filepath.Join(tempDir, ".env"), []byte(dotEnv), 0o600); err != nil {
		t.Fatal(err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	originalValue, wasSet := os.LookupEnv("DOTENV_TEST_FROM_FILE")
	if err := os.Unsetenv("DOTENV_TEST_FROM_FILE"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if wasSet {
			_ = os.Setenv("DOTENV_TEST_FROM_FILE", originalValue)
			return
		}
		_ = os.Unsetenv("DOTENV_TEST_FROM_FILE")
	})
	t.Setenv("DOTENV_TEST_PRECEDENCE", "from-environment")

	if _, err := Load(); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if value := os.Getenv("DOTENV_TEST_FROM_FILE"); value != "loaded" {
		t.Fatalf("DOTENV_TEST_FROM_FILE = %q, want loaded", value)
	}
	if value := os.Getenv("DOTENV_TEST_PRECEDENCE"); value != "from-environment" {
		t.Fatalf("DOTENV_TEST_PRECEDENCE = %q, want existing environment to win", value)
	}
}

func TestLoadAllowsMissingDotEnv(t *testing.T) {
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(originalDir) })

	if _, err := Load(); err != nil {
		t.Fatalf("Load() error = %v, want missing .env to be optional", err)
	}
}
