package controller

import (
	"fmt"
	"github.com/Milefer7/LAN-chat/internal"
	"github.com/Milefer7/LAN-chat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net"
)

// StartServer 开启服务器
func StartServer(r *gin.Engine) {
	log.Println(fmt.Sprintf("Server is running at %s", "localhost:80"))
	if err := r.Run(":80"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

// Ws 处理websocket连接
func Ws(c *gin.Context) {
	// 调用升级器 upgrader 将HTTP请求升级到WebSocket连接将HTTP请求升级到WebSocket连接
	ws, err := utils.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// 关闭连接
	defer func(ws *websocket.Conn) {
		// Send a close frame
		message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		if err := ws.WriteMessage(websocket.CloseMessage, message); err != nil {
			log.Println("Write close message error:", err)
		}
		err := ws.Close()
		if err != nil {
			log.Println("Close error:", err)
		}
	}(ws)

	// 检查客户端是否在本地, 如果是则记录
	remoteAddr := ws.RemoteAddr().String()
	host, _, err := net.SplitHostPort(remoteAddr)
	if host == "127.0.0.1" {
		utils.LocalClient = ws
	}

	// 无限循环处理发来的信息
	internal.HandleGet(ws)
}
