package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
	Storage  StorageConfig  `yaml:"storage"`
}

type AppConfig struct {
	Environment string `yaml:"environment"` // "development", "production"
	DataDir     string `yaml:"data_dir"`
	DeviceName  string `yaml:"device_name"`
	DeviceCode  string `yaml:"device_code"`
}

type DatabaseConfig struct {
	Driver       string `yaml:"driver"` // "sqlite"
	Path         string `yaml:"path"`
	BusyTimeout  int    `yaml:"busy_timeout"`   // in ms, default 5000
	MaxReadConns int    `yaml:"max_read_conns"` // default 4
}

type LoggingConfig struct {
	Level      string `yaml:"level"` // "debug", "info", "warn", "error"
	Format     string `yaml:"format"` // "json", "text"
	FilePath   string `yaml:"file_path"`
	MaxSizeMB  int    `yaml:"max_size_mb"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAgeDays int    `yaml:"max_age_days"`
	Console    bool   `yaml:"console"`
}

type StorageConfig struct {
	DocumentsDir string `yaml:"documents_dir"`
	BackupsDir   string `yaml:"backups_dir"`
}

// DefaultDataDir returns the installation directory (where the application executable is located)
// so the database and runtime data live in the same installed folder, or falls back to
// the local app data directory if the executable folder is read-only.
func DefaultDataDir() (string, error) {
	// First check if the executable directory can be used (installed application folder)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// Ignore temporary directories created by go test / go run
		lowerDir := strings.ToLower(exeDir)
		if !strings.Contains(lowerDir, "temp") && !strings.Contains(lowerDir, "tmp") {
			// Verify write permission
			testFile := filepath.Join(exeDir, ".write_test")
			if err := os.WriteFile(testFile, []byte("ok"), 0o600); err == nil {
				_ = os.Remove(testFile)
				return exeDir, nil
			}
		}
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".school-management"), nil
	}
	return filepath.Join(localAppData, "SchoolManagement"), nil
}

// DefaultConfig generates the default configuration located in baseDir
func DefaultConfig(baseDir string) Config {
	return Config{
		App: AppConfig{
			Environment: "production",
			DataDir:     baseDir,
			DeviceName:  "Primary Desktop",
			DeviceCode:  "DSK1",
		},
		Database: DatabaseConfig{
			Driver:       "sqlite",
			Path:         filepath.Join(baseDir, "data", "school.db"),
			BusyTimeout:  5000,
			MaxReadConns: 4,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			FilePath:   filepath.Join(baseDir, "logs", "app.log"),
			MaxSizeMB:  20,
			MaxBackups: 5,
			MaxAgeDays: 30,
			Console:    true,
		},
		Storage: StorageConfig{
			DocumentsDir: filepath.Join(baseDir, "documents"),
			BackupsDir:   filepath.Join(baseDir, "backups"),
		},
	}
}

// Load loads configuration from configPath, creating default files/dirs if not present.
func Load(configPath string) (*Config, error) {
	var baseDir string
	if configPath == "" {
		var err error
		baseDir, err = DefaultDataDir()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve default data directory: %w", err)
		}
		configPath = filepath.Join(baseDir, "config.yaml")
	} else {
		baseDir = filepath.Dir(configPath)
	}

	cfg := DefaultConfig(baseDir)

	data, err := os.ReadFile(configPath)
	if err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	} else if os.IsNotExist(err) {
		// Create base directories and save default config
		if err := EnsureDirectories(&cfg); err != nil {
			return nil, err
		}
		data, err := yaml.Marshal(&cfg)
		if err == nil {
			_ = os.WriteFile(configPath, data, 0o600)
		}
	} else {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Apply Environment Variable Overrides
	if env := os.Getenv("SMS_ENV"); env != "" {
		cfg.App.Environment = env
	}
	if dbPath := os.Getenv("SMS_DB_PATH"); dbPath != "" {
		cfg.Database.Path = dbPath
	}
	if logLvl := os.Getenv("SMS_LOG_LEVEL"); logLvl != "" {
		cfg.Logging.Level = logLvl
	}

	if err := EnsureDirectories(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// EnsureDirectories ensures that required subdirectories exist.
func EnsureDirectories(cfg *Config) error {
	dirs := []string{
		cfg.App.DataDir,
		filepath.Dir(cfg.Database.Path),
		filepath.Dir(cfg.Logging.FilePath),
		cfg.Storage.DocumentsDir,
		cfg.Storage.BackupsDir,
	}
	for _, dir := range dirs {
		if dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
	}
	return nil
}

