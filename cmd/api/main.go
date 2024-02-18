package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Server struct {
	config Config
}

type Configuration struct {
	Port int
}

func main() {

	var cfg Configuration

	pflag.IntVar(&cfg.Port, "port", 8080, "HTTP port")
	pflag.Parse()

	app := application{
		cfg: cfg,
	s := Server{
		config: cfg,
	}

	err := s.serve()
	if err != nil {
		fmt.Println("server exited with error", err)
	}
}
