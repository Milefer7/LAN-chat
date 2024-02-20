package internal

// 放一些处理消息的函数
import (
	"encoding/json"

	"github.com/Milefer7/LAN-chat/app/model"
	"github.com/Milefer7/LAN-chat/utils"
	"github.com/gorilla/websocket"
	"log"
)

// HandleGet 处理接收到的消息(读取消息 发送响应 调函数将消息展示给用户)
func HandleGet(ws *websocket.Conn) {
	// 无限循环监听
	for {
		// 读取客户端发送的消息
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("读取 error:", err)
			SendFailMsg(ws, model.RequestMessage{})
			continue // 继续监听下一条消息
		}

		// 解析客户端发送的消息
		var reqMsg model.RequestMessage
		if err := json.Unmarshal(message, &reqMsg); err != nil {
			log.Println("解析 error:", err)
			SendFailMsg(ws, reqMsg)
			continue // 继续监听下一条消息
		}
		// 将resMsg发送给qt客户端
		SendMsgToClient(message, ws)
		// 如果成功 将响应消息返回给发送者
		SendSuccessMsg(ws, reqMsg)
		// 调函数将消息展示给用户
		//view.DisplayReq(reqMsg)
	}
}

// SendSuccessMsg 返回成功消息
func SendSuccessMsg(ws *websocket.Conn, reqMsg model.RequestMessage) {
	var respMsg = model.ResponseMessage{
		Data: model.Data{
			Type:      "response",
			Timestamp: reqMsg.Data.Timestamp,
			From:      reqMsg.Data.To,   // 响应方标识设置为请求的接收方标识
			To:        reqMsg.Data.From, // 请求方标识设置为响应的接收方标识
		},
		ContentRes: model.ContentRes{
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
func SendFailMsg(ws *websocket.Conn, reqMsg model.RequestMessage) {
	var respMsg = model.ResponseMessage{
		Data: model.Data{
			Type:      "response",
			Timestamp: reqMsg.Data.Timestamp,
			From:      reqMsg.Data.To,   // 响应方标识设置为请求的接收方标识
			To:        reqMsg.Data.From, // 请求方标识设置为响应的接收方标识
		},
		ContentRes: model.ContentRes{
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
	LocalWs := utils.LocalClient
	// 发送消息
	err := LocalWs.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("Error sending message:", err)
		SendFailMsg(ws, model.RequestMessage{})
	}
}
