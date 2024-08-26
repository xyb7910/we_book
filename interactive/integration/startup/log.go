package startup

import "we_book/pkg/logger"

func InitLog() logger.V1 {
	return logger.NewNoLogger()
}
