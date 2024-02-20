package utils

import (
	"github.com/gorilla/websocket"
	"net/http"
)

// LocalClient 记录本地客户端的ws
var LocalClient *websocket.Conn

// Upgrader 定义升级器
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // 读取缓冲区大小
	WriteBufferSize: 1024, // 写入缓冲区大小
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源访问
	},
}
