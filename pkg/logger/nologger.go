package logger

type NoLogger struct {
}

func (n *NoLogger) String(s string, path string) Field {
	//TODO implement me
	panic("implement me")
}

func NewNoLogger() *NoLogger {
	return &NoLogger{}
}

func (n *NoLogger) Debug(msg string, args ...Field) {
}

func (n *NoLogger) Info(msg string, args ...Field) {
}

func (n *NoLogger) Warn(msg string, args ...Field) {
}

func (n *NoLogger) Error(msg string, args ...Field) {
}
