package configure

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

// не хочу заморачиваться со считыванием из env
var (
	configPath = "./config.yml"
)

type Config struct {
	Server     string        `yaml:"server"`
	AccessTTL  time.Duration `yaml:"accessTTL"`
	RefreshTTL time.Duration `yaml:"refreshTTL"`
	Postgres   Postgres      `yaml:"postgres"`
}

func New() *Config {
	return &Config{
		Server:   "",
		Postgres: Postgres{},
	}
}

func MustConfig() *Config {
	if _, ok := os.Stat(configPath); os.IsNotExist(ok) {
		panic("Config file does not exist: " + configPath)
	}

	cfg := New()

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return cfg
}
