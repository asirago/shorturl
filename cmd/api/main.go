package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Server struct {
	cfg    Config
	logger *slog.Logger
}

type Config struct {
	Port        int    `json:"port"`
	Environment string `json:"environment"`
	LogFile     string `json:"logFile"`
	Output      io.Writer
}

func main() {
	var config Config

	pflag.IntVar(&config.Port, "port", 8080, "HTTP port")
	pflag.StringVar(&config.Environment, "environment", "development", "development|production")
	pflag.StringVar(&config.LogFile, "logFile", "", "Log to File")
	pflag.Parse()

	setupConfigFile("config", &config)

	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		mw := io.MultiWriter(os.Stdout, f)
		config.Output = mw
		defer f.Close()
	}

	s := Server{
		cfg:    config,
		logger: NewLogger(config.Output),
	}

	err := s.serve()
	if err != nil {
		fmt.Println("server exited with error", err)
	}
}

func setupConfigFile(filename string, cfg *Config) error {
	viper.SetConfigName(filename)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	return nil
}
