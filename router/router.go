package router

import (
	"context"

	"github.com/Milefer7/LAN-chat/app/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(e *gin.Engine, cancel context.CancelFunc, ctx context.Context) {
	// 初始化连接接口
	e.GET("/ws", func(c *gin.Context) {
		controller.Ws(c, ctx)
	})
	// 广播接口
	e.POST("/broadcast", func(c *gin.Context) {
		controller.Broadcast(c)
	})
	// 获取在线用户接口
	e.GET("/getOnlineUsers", controller.GetOnlineUsers)
	// 离开接口
	e.POST("/leave", func(c *gin.Context) {
		controller.Leave(c, cancel)
	})
}
