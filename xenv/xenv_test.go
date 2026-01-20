package xenv

import (
	"os"
	"testing"
	"time"
)

type Database struct {
	Host     string `env:"DB_HOST" default:"localhost"`
	Port     int    `env:"DB_PORT" default:"5432"`
	User     string `env:"DB_USER" required:"true"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	SSLMode  bool   `env:"DB_SSLMODE" default:"false"`
}

type AppConfig struct {
	Debug   bool          `env:"DEBUG" default:"false"`
	Timeout time.Duration `env:"TIMEOUT" default:"30s"`
	APIKey  string        `env:"KEY" required:"true"`
	Metrics struct {
		Enabled bool `env:"ENABLED" default:"true"`
	}
	Database Database `env:"DB"`
}

func TestBasicConfig(t *testing.T) {
	os.Setenv("KEY", "secret-123")
	os.Setenv("DEBUG", "true")
	// set db envs
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "dbpass")
	os.Setenv("DB_SSLMODE", "true")

	cfg := &AppConfig{}
	err := Load(cfg)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Debug != true {
		t.Errorf("Expected true, got %v", cfg.Debug)
	}

	if cfg.APIKey != "secret-123" {
		t.Errorf("Expected 'secret-123', got %s", cfg.APIKey)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected 'localhost', got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected 5432, got %d", cfg.Database.Port)
	}

	if cfg.Database.User != "user" {
		t.Errorf("Expected 'user', got %s", cfg.Database.User)
	}

	if cfg.Database.Password != "dbpass" {
		t.Errorf("Expected 'dbpass', got %s", cfg.Database.Password)
	}

	if !cfg.Database.SSLMode {
		t.Errorf("Expected 'true', got %t", cfg.Database.SSLMode)
	}

	defer os.Clearenv()
}

func TestAdvancedFeatures(t *testing.T) {
	// We set vars with a prefix "APP_"
	os.Setenv("APP_KEY", "secret-123")
	os.Setenv("APP_TIMEOUT", "5m")
	os.Setenv("DEBUG", "true")
	defer os.Clearenv()

	cfg := &AppConfig{}
	opts := Options{Prefix: "APP_"}

	err := LoadWithOptions(cfg, opts)

	if err == nil {
		t.Fatalf("Load should have failed because `APP_DB_USER` us required: %v", err)
	}

	os.Setenv("APP_DB_USER", "user")
	err = LoadWithOptions(cfg, opts)

	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Test Debug not included
	if cfg.Debug != false {
		t.Errorf("Expected false, got %v", cfg.Debug)
	}

	// Test Duration
	if cfg.Timeout != 5*time.Minute {
		t.Errorf("Expected 5m, got %v", cfg.Timeout)
	}

	// Test Prefix Mapping
	if cfg.APIKey != "secret-123" {
		t.Errorf("Prefix mapping failed, got %s", cfg.APIKey)
	}

	// Test Nested with Prefix
	if cfg.Metrics.Enabled != true {
		t.Error("Nested default failed")
	}
}
