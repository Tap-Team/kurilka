package app

import (
	"os"

	"log/slog"
)

func SetLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
