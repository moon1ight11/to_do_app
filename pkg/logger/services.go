package logger

import "os"

// закрытие файла логов
func (l *Logger) Close() error {
	return l.file.Close()
}

// логирование информационного сообщения
func (l *Logger) Info(msg string, fields ...any) {
	l.log("INFO", msg, fields...)
}

// логирование сообщения об ошибке
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields...)
}

// логирование фатальной ошбки с завершением программы
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}
