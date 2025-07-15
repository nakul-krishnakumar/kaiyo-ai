package config

import (
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

func MustLoad() *Config {
	// load env files
	if err := godotenv.Load(); err != nil {
        slog.Warn("No .env file found or error loading it", slog.String("error", err.Error()))
    }

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			slog.Error("Config path is not set")
			os.Exit(1)
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		slog.Error("Config file does not exist", slog.String("config-path", configPath))
		os.Exit(1)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		slog.Error("Failed to load config file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// set port from cmd
	port := flag.String("port", "", "port to expose the app")
	flag.Parse()
	if *port != "" {
		cfg.HTTPServer.Port = *port
	}

	return &cfg
}
