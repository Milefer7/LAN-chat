package broadcast

import (
	"sync"
	"time"
)

// 定义GetIn接口的请求消息结构

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

//type MsgBack struct {
//	Code    string `json:"code"`
//	Message string `json:"message"`
//	Users   []User `json:"data"`
//}

// 设定已有在线用户列表和互斥锁，用于线程安全地更新和访问用户列表
var OnlineUsers []User
var OnlineUsersMutex sync.Mutex

// 超时时间为30秒
const HeartbeatTimeout = 30 * time.Second
