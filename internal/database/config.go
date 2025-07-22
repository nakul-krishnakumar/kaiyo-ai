package database

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	URL  string     
	Pool PoolConfig 
}

type PoolConfig struct {
	URL         string
	MinConns    int32         `mapstructure:"min_conns"`
	MaxConns    int32         `mapstructure:"max_conns"`
	MaxConnLife time.Duration `mapstructure:"max_conn_life"`
	MaxConnIdle time.Duration `mapstructure:"max_conn_idle"`

	/*
		URL -- connection string
		MinConns -- min simultaneous conns
		MaxConns -- max simultaneous conns
		MaxConnLife -- to refresh each conn to maintain good conns
		MaxConnIdle -- max idle time for each conn
	*/
}

type TimeoutConfig struct {
	Connect time.Duration `mapstructure:"connect"`
	Query   time.Duration `mapstructure:"query"`
}

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Error("Unable to load database secrets" + err.Error())
	}

	var cfg Config

	url := os.Getenv("DATABASE_URL")
	cfg.URL = url

	configPath := os.Getenv("DB_CONFIG_PATH")
	if configPath == "" {
		configPath = "/config/database.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		slog.Error("Failed to read database config file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		slog.Error("Failed to load database config file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return &cfg
}
