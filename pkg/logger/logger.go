package logger

import (
	"log/slog"
	"os"
	"sso/pkg/logger/handler"
)

func New(f *os.File) *slog.Logger {
	logger := slog.New(handler.NewCopyHandler(slog.NewTextHandler(os.Stdout, nil), slog.NewJSONHandler(f, nil)))
	return logger
}
