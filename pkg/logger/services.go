package logger

import "os"

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) Info(msg string, fields ...any) {
	l.log("INFO", msg, fields...)
}

func (l *Logger) Error(msg string, fields ...any) {
	l.log("ERROR", msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...any) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}
