package main

import (
	"github.com/gin-gonic/gin"
	corn "github.com/robfig/cron/v3"
	"we_book/events"
)

type App struct {
	web      *gin.Engine
	consumer []events.Consumer
	corn     *corn.Cron
}
