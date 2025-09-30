package logger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var FilePathNotSet error = errors.New("LOG_FILE_PATH environment variable not set")

type Logger struct {
	info     *slog.Logger
	error    *slog.Logger
	file     *slog.Logger
	FileOnly bool
}

func NewSafeLogger(fileOnly bool) (*Logger, error) {
	logFilePath := os.Getenv("LOG_FILE_PATH")
	if logFilePath == "" {
		return nil, FilePathNotSet
	}
	filePath := fmt.Sprintf("app-%v.log", time.Now().UnixNano())
	absPath, err := filepath.Abs(logFilePath + filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path of log file: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(absPath), 0o755); err != nil {
		return nil, fmt.Errorf("failed to make log directory: %w", err)
	}
	logFile, err := os.OpenFile(logFilePath+filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return &Logger{
		info: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
		error: slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelError,
		})),
		file: slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
		FileOnly: fileOnly,
	}, nil
}

func NewLogger(fileOnly bool) *Logger {
	logger, err := NewSafeLogger(fileOnly)
	if err != nil {
		panic(err)
	}
	return logger
}

func (l *Logger) Debug(msg string, args ...any) {
	if !l.FileOnly {
		l.info.Debug(msg, args...)
	}
	if os.Getenv("DEBUG") == "1" {
		l.file.Debug(msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...any) {
	if !l.FileOnly {
		l.info.Info(msg, args...)
	}
	l.file.Info(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	if !l.FileOnly {
		l.info.Warn(msg, args...)
	}
	l.file.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	if !l.FileOnly {
		l.error.Error(msg, args...)
	}
	l.file.Error(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	if !l.FileOnly {
		l.error.Error(msg, args...)
	}
	l.file.Error(msg, args...)
	os.Exit(1)
}
