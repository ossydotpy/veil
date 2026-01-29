package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	MasterKey string
	DbPath    string
	StoreType string
}

func (c *Config) Validate() error {
	if c.DbPath == "" {
		return errors.New("database path is required")
	}
	if c.StoreType == "" {
		return errors.New("STORE_TYPE is required")
	}
	return nil
}

func (c *Config) ValidateMasterKey() error {
	if c.MasterKey == "" {
		return errors.New("MASTER_KEY is required")
	}
	if len(c.MasterKey) != 64 {
		return fmt.Errorf("invalid MASTER_KEY: expected 64 hex characters (32 bytes), got %d", len(c.MasterKey))
	}
	return nil
}

func LoadConfig() *Config {
	masterkey := getenv("MASTER_KEY", "")

	home, _ := os.UserHomeDir()
	defaultDbPath := filepath.Join(home, ".veil.db")
	dbPath := getenv("VEIL_DB_PATH", defaultDbPath)
	storeType := getenv("VEIL_STORE_TYPE", "sqlite")

	cfg := &Config{
		MasterKey: masterkey,
		DbPath:    dbPath,
		StoreType: storeType,
	}
	return cfg
}

func getenv(key string, def string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return def
}
