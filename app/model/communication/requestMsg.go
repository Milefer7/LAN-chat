package communication

import "github.com/Milefer7/LAN-chat/app/model/broadcast"

// 定义requestMsg数据结构

type Data struct {
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	From      string `json:"from"`
	To        string `json:"to"`
}

type ContentReq struct {
	PartType string `json:"partType"`
	Data     string `json:"data"`
}

type OnlineUsers struct {
	OnlineUsers []broadcast.User `json:"OnlineUsers"`
}

type RequestMessage struct {
	Data        `json:"data"`
	Content     []ContentReq `json:"content"`
	OnlineUsers `json:"OnlineUsers"`
}
