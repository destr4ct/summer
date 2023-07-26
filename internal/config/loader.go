package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
	"time"
)

var (
	defaultPath = "./etc/default.yml"
	onceHook    = new(sync.Once)
	cfg         *Config
)

func init() {
	if envPath := os.Getenv("CONFIG_PATH"); len(envPath) > 0 {
		defaultPath = envPath
	}
}

type Config struct {
	Env   string        `yaml:"env" env:"ENV" env-default:"dev"`
	Delay time.Duration `yaml:"delay" env:"DELAY" env-default:"1h"`

	BrokerConfig   BrokerConfig   `yaml:"broker"`
	TelegramConfig TelegramConfig `yaml:"telegram"`
	DBConfig       DatabaseConfig `yaml:"db"`
}

type BrokerConfig struct {
	Host string `yaml:"host" env:"BRK_HOST"`
	Port int    `yaml:"port" env:"BRK_PORT"`

	Username string `env:"BRK_USERNAME"`
	Password string `env:"BRK_PASSWORD"`
}

type DatabaseConfig struct {
	Host string `yaml:"host" env:"DB_HOST"`
	Port int    `yaml:"port" env:"DB_PORT"`

	Username string `env:"DB_USERNAME"`
	Password string `env:"DB_PASSWORD"`
	DB       string `yaml:"db_name" env:"DB_NAME"`
}

type TelegramConfig struct {
	APIKey string `yaml:"api_key" env:"TG_API_KEY"`
}

func Load(path ...string) *Config {
	// конфиг - синглтон
	onceHook.Do(func() {

		// Если был передан путь, то используем его, а не тот, что был передан по умолчанию
		loadPath := defaultPath
		if len(path) == 1 {
			loadPath = path[0]
		}

		cfg = &Config{}
		if err := cleanenv.ReadConfig(loadPath, cfg); err != nil {
			log.Panicf("config.Load: failed to parse config: %v\n", err)
		}

		if err := cleanenv.ReadEnv(cfg); err != nil {
			log.Panicf("config.Load: failed to load env: %v\n", err)
		}

	})
	return cfg
}
