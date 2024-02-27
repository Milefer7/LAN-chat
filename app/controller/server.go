package controller

import "C"
import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Milefer7/LAN-chat/app/model/broadcast"
	"github.com/Milefer7/LAN-chat/internal"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StartServer 开启服务器
func StartServer(r *gin.Engine) {
	log.Println(fmt.Sprintf("Server is running at %s", "localhost:80"))
	if err := r.Run(":80"); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}

// Ws 处理websocket连接
func Ws(c *gin.Context, ctx context.Context) {
	// 调用升级器 upgrader 将HTTP请求升级到WebSocket连接将HTTP请求升级到WebSocket连接
	ws, err := broadcast.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// 关闭连接
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Println("Close error:", err)
		}
	}(ws)

	// 检查客户端是否在本地, 如果是则记录
	remoteAddr := ws.RemoteAddr().String()
	host, _, err := net.SplitHostPort(remoteAddr)
	if host == "127.0.0.1" {
		broadcast.LocalClient = ws
		log.Println("本地客户端已连接")
	}

	// 无限循环处理客户端发来的信息
	internal.HandleGet(ws, ctx)
}

// Broadcast 广播 join消息+返回在线用户
func Broadcast(c *gin.Context) {
	// 解析body
	flag := true
	var data broadcast.BroadcastMsg
	if err := c.BindJSON(&data); err != nil {
		log.Println("发送广播 BindJSON error:", err)
		flag = false
	}
	// 发送广播消息
	err := internal.CreateUdpConn(&data)
	if err != nil {
		log.Println("发送广播 CreateUdpConn error:", err)
		flag = false
	}
	// 返回在线用户
	if flag {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "发送广播成功",
		})
	} else {
		c.JSON(200, gin.H{
			"code":    0,
			"message": "调用广播接口失败，详细信息请查看日志",
		})
	}
}

// Leave 定义用户离开接口
func Leave(c *gin.Context, cancel context.CancelFunc) {
	var req broadcast.LeaveRequest
	if err := c.BindJSON(&req); err != nil {
		log.Println("离开接口 BindJSON error:", err)
		return
	}
	if req.Leave {
		broadcast.LocalClient = nil
		// 通知goroutine退出
		defer cancel()
	}
	c.JSON(200, gin.H{
		"code":    "1",
		"message": "调用离开接口成功",
	})
}

func GetOnlineUsers(c *gin.Context) {
	// 返回在线用户
	c.JSON(200, gin.H{
		"code":        1,
		"message":     "获取在线用户成功",
		"onlineUsers": broadcast.OnlineUsers,
	})
}
