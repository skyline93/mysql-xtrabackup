package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Identifer   string
	Version     string
	LoginPath   string
	DbHostName  string
	DbUser      string
	Throttle    int
	TryCompress bool

	BinPath        string
	DataPath       string
	BackupUser     string
	BackupHostName string
}

func saveConfigToRepo(config *Config, repoPath string) error {
	d, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(repoPath, "config"), d, 0664); err != nil {
		return err
	}

	return nil
}

func loadConfigFromRepo(config *Config, path string) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(d, config); err != nil {
		return err
	}

	return nil
}
