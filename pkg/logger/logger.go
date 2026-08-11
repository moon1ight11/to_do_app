package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"todoapp/internal/config"
)

type Logger struct {
	file   *os.File
	logger *log.Logger
}

func New(cfg *config.Config) (*Logger, error) {
	dir := filepath.Dir(cfg.Logger.FilePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("logger.New: create log directory: %w", err)
		}
	}

	file, err := os.OpenFile(cfg.Logger.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("logger.New: open log file: %w", err)
	}

	return &Logger{
		file:   file,
		logger: log.New(file, "", 0),
	}, nil
}

func (l *Logger) log(level string, msg string, fields ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var b strings.Builder
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			if b.Len() > 0 {
				b.WriteString(" ")
			}
			b.WriteString(fmt.Sprintf("%v=%v", fields[i], fields[i+1]))
		}
	}

	fieldsStr := b.String()
	if fieldsStr != "" {
		fieldsStr = " " + fieldsStr
	}

	l.logger.Printf("[%s] %s: %s%s", timestamp, level, msg, fieldsStr)
}
