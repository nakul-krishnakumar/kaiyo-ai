package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type HTTPServer struct {
	Host string 
	Port string	
}

type Config struct {
	Env string `mapstructure:"env"`
	HTTPServer `mapstructure:"http_server"`
}

func MustLoad(env string, path string) *Config {
	v := viper.New()
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", slog.String("error", err.Error()))
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		slog.Error("Failed to load config file", slog.String("error", err.Error()))
	}

	return &cfg
}
