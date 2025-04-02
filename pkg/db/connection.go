package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

// Config verilənlər bazası konfiqurasiyasını saxlayır
type Config struct {
	ConnectionString string `yaml:"connection_string"`
}

// Connect verilənlər bazasına bağlantı yaradır
func Connect() (*sqlx.DB, error) {
	// Konfiqurasiya faylının oxunması
	configPath := filepath.Join("configs", "db.yaml")
	data, err := os.ReadFile(configPath) // ioutil.ReadFile əvəzinə os.ReadFile
	if err != nil {
		return nil, fmt.Errorf("db konfiqurasiyasının oxunması xətası: %w", err)
	}

	// Konfiqurasiya emalı
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("db konfigurasiyasının emalı xətası: %w", err)
	}

	// Verilənlər bazasına qoşulma
	db, err := sqlx.Connect("postgres", config.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("verilənlər bazasına qoşulma xətası: %w", err)
	}

	// Bağlantının yoxlanması
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("verilənlər bazasına ping xətası: %w", err)
	}

	return db, nil
}
