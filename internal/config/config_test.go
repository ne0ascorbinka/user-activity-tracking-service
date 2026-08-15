package config

import (
	"os"
	"testing"
	"time"
)

func TestConfigLoadAndDSN(t *testing.T) {
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("DB_MAX_CONNS", "50")
	os.Setenv("DB_MIN_CONNS", "10")
	os.Setenv("DB_MAX_CONN_LIFETIME", "2h")
	os.Setenv("DB_MAX_CONN_IDLE_TIME", "45m")
	os.Setenv("AGGREGATION_INTERVAL", "2h")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if cfg.ServerPort != 9000 {
		t.Errorf("expected ServerPort=9000, got %d", cfg.ServerPort)
	}
	if cfg.DBHost != "db.example.com" {
		t.Errorf("expected DBHost=db.example.com, got %s", cfg.DBHost)
	}
	if cfg.DBPort != 5433 {
		t.Errorf("expected DBPort=5433, got %d", cfg.DBPort)
	}
	if cfg.DBUser != "testuser" {
		t.Errorf("expected DBUser=testuser, got %s", cfg.DBUser)
	}
	if cfg.DBPassword != "secret" {
		t.Errorf("expected DBPassword=secret, got %s", cfg.DBPassword)
	}
	if cfg.DBName != "testdb" {
		t.Errorf("expected DBName=testdb, got %s", cfg.DBName)
	}
	if cfg.DBSSLMode != "require" {
		t.Errorf("expected DBSSLMode=require, got %s", cfg.DBSSLMode)
	}
	if cfg.DBMaxConns != 50 {
		t.Errorf("expected DBMaxConns=50, got %d", cfg.DBMaxConns)
	}
	if cfg.DBMinConns != 10 {
		t.Errorf("expected DBMinConns=10, got %d", cfg.DBMinConns)
	}
	if cfg.DBMaxConnLifetime != 2*time.Hour {
		t.Errorf("expected DBMaxConnLifetime=2h, got %v", cfg.DBMaxConnLifetime)
	}
	if cfg.DBMaxConnIdleTime != 45*time.Minute {
		t.Errorf("expected DBMaxConnIdleTime=45m, got %v", cfg.DBMaxConnIdleTime)
	}
	if cfg.AggregationInterval != 2*time.Hour {
		t.Errorf("expected AggregationInterval=2h, got %v", cfg.AggregationInterval)
	}

	expectedDSN := "postgres://testuser:secret@db.example.com:5433/testdb?sslmode=require"
	if cfg.DSN() != expectedDSN {
		t.Errorf("expected DSN=%s, got %s", expectedDSN, cfg.DSN())
	}
}
