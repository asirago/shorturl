package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func (s *Server) serve() error {

	var addr string = fmt.Sprintf("localhost:%d", s.cfg.Port)
	if s.cfg.Environment == "production" {
		addr = fmt.Sprintf(":%v", strings.Split(addr, ":")[1])
	}
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		fmt.Printf("caught signal: %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// complete background tasks here with waitgroup

		shutdownError <- nil

	}()

	fmt.Printf("starting server on port %v\n", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return nil
	}

	fmt.Printf("stoped server on port %v\n", srv.Addr)

	return nil

}
