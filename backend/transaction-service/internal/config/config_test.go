package config

import "testing"

func TestLoadDefaultsToDedicatedLocalPort(t *testing.T) {
	t.Setenv("SERVICE_PORT", "")
	t.Setenv("BOOK_SERVICE_URL", "")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config.Service.Port != ":8084" {
		t.Fatalf("Load() port = %q, want :8084", config.Service.Port)
	}
	if config.Book.URL != "http://localhost:8083" {
		t.Fatalf("Load() book URL = %q, want http://localhost:8083", config.Book.URL)
	}
}
