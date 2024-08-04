package logger

type V1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	String(s string, path string) Field
}

type Field struct {
	Key   string
	Value any
}
