package logger

import (
	"log/slog"
	"os"
)

var (
	Log *slog.Logger
)

func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// Switched to TextHandler to prevent PowerShell terminal mangling from explicit \n string breaks inside raw JSON error traces
	handler := slog.NewTextHandler(os.Stdout, opts)

	Log = slog.New(handler)
	slog.SetDefault(Log) // Globally sets standard logger allowing implicit mappings
}
