package internal

// 放一些处理消息的函数
import (
	"context"
	"encoding/json"
	"github.com/Milefer7/LAN-chat/app/model/broadcast"
	"github.com/Milefer7/LAN-chat/app/model/communication"
	"time"

	"github.com/gorilla/websocket"
	"log"
)

// HandleGet 处理接收到的消息(读取消息 发送响应 调函数将消息展示给用户)
func HandleGet(ws *websocket.Conn, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// 当接收到取消信号时，停止循环
			log.Println("接受到取消信号，停止处理消息")
			return
		default:
			// 读取客户端发送的消息
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("读取 error:", err)
				SendFailMsg(ws, communication.RequestMessage{})
				continue // 继续监听下一条消息
			}

			// 解析客户端发送的消息
			var reqMsg communication.RequestMessage
			if err := json.Unmarshal(message, &reqMsg); err != nil {
				log.Println("解析客户端发送的消息错误 error:", err)
				SendFailMsg(ws, reqMsg)
				continue // 继续监听下一条消息
			}
			// 处理消息
			// 将resMsg发送给qt客户端
			SendMsgToClient(message, ws)
			// 如果成功 将响应消息返回给发送者
			SendSuccessMsg(ws, reqMsg)
		}
	}
}

// SendSuccessMsg 返回成功消息
func SendSuccessMsg(ws *websocket.Conn, reqMsg communication.RequestMessage) {
	var respMsg = communication.ResponseMessage{
		Data: communication.Data{
			Type:      "response",
			Timestamp: reqMsg.Data.Timestamp,
			From:      reqMsg.Data.To,   // 响应方标识设置为请求的接收方标识
			To:        reqMsg.Data.From, // 请求方标识设置为响应的接收方标识
		},
		ContentRes: communication.ContentRes{
			Status:  1,
			Message: "消息接收成功",
		},
	}

	// 将respMsg转为json格式
	resp, err := json.Marshal(respMsg)
	if err != nil {
		log.Println("Marshal error:", err)
	}

	if err := ws.WriteMessage(websocket.TextMessage, resp); err != nil {
		log.Println("Write error:", err)
	}
}

// SendFailMsg 返回失败消息
func SendFailMsg(ws *websocket.Conn, reqMsg communication.RequestMessage) {
	var respMsg = communication.ResponseMessage{
		Data: communication.Data{
			Type:      "response",
			Timestamp: reqMsg.Data.Timestamp,
			From:      reqMsg.Data.To,   // 响应方标识设置为请求的接收方标识
			To:        reqMsg.Data.From, // 请求方标识设置为响应的接收方标识
		},
		ContentRes: communication.ContentRes{
			Status:  0,
			Message: "消息接收失败",
		},
	}

	// 将respMsg转为json格式
	resp, err := json.Marshal(respMsg)
	if err != nil {
		log.Println("Marshal error:", err)
	}

	if err := ws.WriteMessage(websocket.TextMessage, resp); err != nil {
		log.Println("Write error:", err)
	}
}

// SendMsgToClient 给本地客户端发送消息, 传入的ws是局域网用户的ws
func SendMsgToClient(message []byte, ws *websocket.Conn) {
	// 获取客户端ws
	LocalWs := LocalClient
	// 发送消息
	err := LocalWs.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Error sending message:", err)
		SendFailMsg(ws, communication.RequestMessage{})
	}
}

// SendOnlineUsers 返回在线用户
func SendOnlineUsers(ws *websocket.Conn, ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var requestMsg = communication.RequestMessage{
				Data:    communication.Data{},
				Content: []communication.ContentReq{},
				OnlineUsers: communication.OnlineUsers{
					OnlineUsers: broadcast.OnlineUsers,
				},
			}
			// 将OnlineUsers转为json格式
			users, err := json.Marshal(requestMsg)
			if err != nil {
				log.Println("Marshal error:", err)
				continue
			}
			// 发送在线用户
			if err := ws.WriteMessage(websocket.TextMessage, users); err != nil {
				log.Println("Write error:", err)
			}
		case <-ctx.Done():
			log.Println("停止发送在线用户")
			// 停止
			return
		}
	}
}
