package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// 确定 logger 的编码格式
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

// 将日志写到哪里
func getWriteSyncer() zapcore.WriteSyncer {
	file, _ := os.Create("./logs/app.log")
	return zapcore.AddSync(file)
}

// InitSelfLogger 初始化 logger
func InitSelfLogger() *zap.SugaredLogger {
	writeSyncer := getWriteSyncer()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	SugLogger := logger.Sugar()
	return SugLogger
}
