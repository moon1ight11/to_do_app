package logger

type LoggerInterface interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	Close() error
}
