package slogx

import (
	"io"
	"log/slog"
	"os"
)

type logger struct {
	handler slog.Handler
	w       io.Writer
	level   slog.Level
	source  bool
}

func New(options ...Option) *slog.Logger {
	var l logger
	for _, opt := range options {
		opt(&l)
	}

	if l.w == nil {
		l.w = os.Stdout
	}

	if l.handler == nil {
		l.handler = slog.NewJSONHandler(l.w, &slog.HandlerOptions{
			Level:     l.level,
			AddSource: l.source,
		})
	}

	return slog.New(l.handler)
}

type Option func(s *logger)

func WithSource(source bool) Option {
	return func(s *logger) {
		s.source = source
	}
}

func WithLevel(level slog.Level) Option {
	return func(s *logger) {
		s.level = level
	}
}

func WithWriter(w io.Writer) Option {
	return func(s *logger) {
		s.w = w
	}
}
