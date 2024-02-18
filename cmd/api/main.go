package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
)

type Server struct {
	config Config
	logger *slog.Logger
}

type Config struct {
	Port        int
	LogToFile   string
	Environment string
	Output      io.Writer
}

func main() {

	var cfg Config

	cfg.Output = os.Stdout

	pflag.IntVar(&cfg.Port, "port", 8080, "HTTP port")
	pflag.StringVar(&cfg.Environment, "environment", "development", "development|production")
	pflag.StringVar(&cfg.LogToFile, "logFile", "", "Log to File")
	pflag.Parse()

	if cfg.LogToFile != "" {
		f, err := os.OpenFile(cfg.LogToFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		mw := io.MultiWriter(os.Stdout, f)
		cfg.Output = mw
		defer f.Close()
	}

	s := Server{
		logger: NewLogger(cfg.Output),
		config: cfg,
	}

	err := s.serve()
	if err != nil {
		fmt.Println("server exited with error", err)
	}
}
