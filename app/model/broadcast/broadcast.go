package broadcast

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// 定义在线用户列表和互斥锁
var OnlineUsers []User
var OnlineUsersMutex sync.Mutex

// 定义超时时间 30秒
const HeartbeatTimeout = 30 * time.Second

// LocalClient 记录本地客户端的ws
var LocalClient *websocket.Conn

//// MultiBroadcastConn 记录多播连接
//var MultiBroadcastConn *net.UDPConn

// LocalBroadcastMsg 记录本地广播消息
var LocalBroadcastMsg BroadcastMsg

// StartHeartbeat 用于启动心跳
var StartHeartbeat = make(chan bool)

// Upgrader 定义升级器
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // 读取缓冲区大小
	WriteBufferSize: 1024, // 写入缓冲区大小
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源访问
	},
}

// getIn接口，前端给后端的消息，也是后端广播的消息
type BroadcastMsg struct {
	Type        string `json:"type"` // 消息类型: join, heartbeat, state_request, state_response
	UserName    string `json:"userName"`
	DeviceType  string `json:"deviceType"`
	Fingerprint string `json:"fingerprint"`
	Port        int    `json:"port"`
}

// 用户信息（返回getIn接口body中的一部分）
type User struct {
	UserName      string    `json:"userName"`
	Fingerprint   string    `json:"fingerprint"`
	Host          string    `json:"host"`
	LastHeartbeat time.Time `json:"lastHeartbeat"`
}

type LeaveRequest struct {
	Leave   bool   `json:"leave"`
	Message string `json:"message"`
}
