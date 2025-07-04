package pkg

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigLoader interface {
	Load(path string, cfg any) error
}

type DBConfig interface {
	DriverName() string
	DSN() string
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	// sqlite
	Path string `yaml:"path,omitempty"`
	// postgres
	Host     string `yaml:"host,omitempty"`
	Port     string `yaml:"port,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Name     string `yaml:"name,omitempty"`
	SSLMODE  string `yaml:"sslmode,omitempty"`
}

type ServerConfig struct {
	Address     string        `yaml:"server" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

// ------ easy connect ---------
func (d DatabaseConfig) DriverName() string {
	return d.Driver
}

func (d DatabaseConfig) DSN() string {
	switch d.Driver {
	case "sqlite":
		return d.Path
	case "postgres":
		return "host=" + d.Host +
			" port=" + d.Port +
			" user=" + d.User +
			" password=" + d.Password +
			" dbname=" + d.Name +
			" sslmode=" + d.SSLMODE
	default:
		return ""
	}
}

// --------- CleanENV ------------

type CleanenvLoader struct{}

func (l CleanenvLoader) Load(path string, cfg any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	expanded := os.ExpandEnv(string(content))

	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(expanded); err != nil {
		return err
	}
	tmpFile.Close()

	return cleanenv.ReadConfig(tmpFile.Name(), cfg)
}


// ------ init -------

func InitConfig(loader ConfigLoader, path string) *Config {
	var cfg Config

	if err := loader.Load(path, &cfg); err != nil {
		panic("ошибка загрузки конфигурации: " + err.Error())
	}

	return &cfg
}