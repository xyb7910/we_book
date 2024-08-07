package ioc

import (
	"go.uber.org/zap"
	logger2 "we_book/pkg/logger"
)

func InitLogger() logger2.V1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger2.NewZapLogger(l)
}
