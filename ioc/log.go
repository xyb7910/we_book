package ioc

import (
	"go.uber.org/zap"
	"we_book/internal/pkg/logger"
)

func InitLogger() logger.V1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
