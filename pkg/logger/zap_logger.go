package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{l: l}
}

func (zl *ZapLogger) toZapFields(args []Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}

func (zl *ZapLogger) Debug(msg string, args ...Field) {
	zl.l.Debug(msg, zl.toZapFields(args)...)
}
func (zl *ZapLogger) Info(msg string, args ...Field) {
	zl.l.Info(msg, zl.toZapFields(args)...)
}
func (zl *ZapLogger) Warn(msg string, args ...Field) {
	zl.l.Warn(msg, zl.toZapFields(args)...)
}
func (zl *ZapLogger) Error(msg string, args ...Field) {
	zl.l.Error(msg, zl.toZapFields(args)...)
}

func (zl *ZapLogger) String(s string, path string) Field {
	panic("implement me")
}
