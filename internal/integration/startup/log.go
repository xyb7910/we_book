package startup

import "we_book/internal/pkg/logger"

//func InitLogger() logger.V1 {
//	l, err := zap.NewDevelopment()
//	if err != nil {
//		panic(err)
//	}
//	return logger.NewZapLogger(l)
//}

func InitLogger() logger.V1 {
	return logger.NewNoLogger()
}
