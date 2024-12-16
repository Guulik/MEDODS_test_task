package configure

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string   `yaml:"env"`
	Server   string   `yaml:"server"`
	Auth     Auth     `yaml:"auth"`
	Postgres Postgres `yaml:"postgres"`
	SMTP     SMTP     `yaml:"smtp"`
}

type Auth struct {
	AccessTTL  time.Duration `yaml:"accessTTL"`
	RefreshTTL time.Duration `yaml:"refreshTTL"`
}

type SMTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func New() *Config {
	return &Config{
		Server:   "",
		Postgres: Postgres{},
	}
}

func MustConfig() *Config {
	configPath := fetchConfigPath()

	if _, ok := os.Stat(configPath); os.IsNotExist(ok) {
		panic("Config file does not exist: " + configPath)
	}

	cfg := New()

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	return "./local.yml"
}
