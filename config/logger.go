package config

import (
	"log/slog"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *slog.Logger

func InitLogger() *slog.Logger {
	rotator := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	handler := slog.NewJSONHandler(rotator, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Logger = slog.New(handler)
	slog.SetDefault(Logger)

	return Logger
}
