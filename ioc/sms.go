package ioc

import (
	"we_book/internal/service/sms"
	"we_book/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
