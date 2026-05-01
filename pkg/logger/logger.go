package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"todoapp/internal/config"
)

// структура логгера
type Logger struct {
	file   *os.File
	logger *log.Logger
}

// конструктор логгера
func New(cfg *config.Config) (*Logger, error) {
	dir := filepath.Dir(cfg.Logger.FilePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
	}

	file, err := os.OpenFile(cfg.Logger.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := log.New(file, "", 0)

	return &Logger{
		file:   file,
		logger: logger,
	}, nil
}

// форматирование и запись
func (l *Logger) log(level string, msg string, fields ...any) {
	// устанавливаем временную метку
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// форматируем сообщения
	fieldsStr := ""
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			fieldsStr += fmt.Sprintf(" %v=%v", fields[i], fields[i+1])
		}
	}

	// записываем в файл
	l.logger.Printf("[%s] %s: %s%s", timestamp, level, msg, fieldsStr)
}
