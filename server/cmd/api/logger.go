package main

import (
	"io"
	"log/slog"
	"time"
)

func NewLogger(out io.Writer) *slog.Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.MessageKey && a.Value.String() == "" {
			return slog.Attr{}
		}
		if a.Key == slog.TimeKey {
			a.Value = slog.AnyValue(time.Now().Format(time.RFC3339))
		}
		return a
	}

	return slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: replace,
	}))
}
