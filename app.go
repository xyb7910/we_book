package main

import (
	"github.com/gin-gonic/gin"
	"we_book/events"
)

type App struct {
	web      *gin.Engine
	consumer []events.Consumer
}
