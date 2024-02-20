package router

import (
	"github.com/Milefer7/LAN-chat/app/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(e *gin.Engine) {
	e.GET("/ws", controller.Ws)
}
