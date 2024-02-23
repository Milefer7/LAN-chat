package controller

import (
	"encoding/json"
	"fmt"
	model "github.com/Milefer7/LAN-chat/app/model/broadcast"
	"github.com/Milefer7/LAN-chat/app/model/communication"
	"log"
	"net"
	"testing"
	"time"
)

func MockBroadcastMsg(i int) model.BroadcastMsg {
	return model.BroadcastMsg{
		UserName:    fmt.Sprintf("TestUser%d", i),
		DeviceType:  fmt.Sprintf("TestDevice%d", i),
		Fingerprint: fmt.Sprintf("Fingerprint%d", i),
		Port:        52017 + i, // 确保端口在合法范围内
	}
}

func TestGetIn(t *testing.T) {
	// 模拟发送广播消息
	conn, err := net.Dial("udp", "224.0.0.167:52017")
	if err != nil {
		log.Fatalf("Dial error: %v", err)
	}
	defer conn.Close()

	// MockBroadcastMsg 生成一个mock的BroadcastMsg

	msg := model.BroadcastMsg{
		UserName:    "TestUser",
		DeviceType:  "TestDevice",
		Fingerprint: "Fingerprint",
		Port:        52017,
	}
	// 打印msg
	fmt.Println(msg)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}

	_, err = conn.Write(msgBytes)
	if err != nil {
		log.Fatalf("Write error: %v", err)
	}
	fmt.Println("广播消息已发送")
}

//	func TestMulticastOnSingleMachine(t *testing.T) {
//		// 创建一个WaitGroup实例
//		var wg sync.WaitGroup
//
//		// 启动多个实例
//		for i := 0; i < 5; i++ {
//			// 每个实例都需要增加一个计数
//			wg.Add(1)
//			go func(i int) {
//				defer wg.Done()
//
//				// 每个实例使用不同的用户名和指纹
//				//userName := fmt.Sprintf("TestUser%d", i)
//				//fingerprint := fmt.Sprintf("Fingerprint%d", i)
//
//				// 创建一个模拟的*gin.Context
//				w := httptest.NewRecorder()
//				c, _ := gin.CreateTestContext(w)
//
//				// 创建一个模拟的BroadcastMsg
//				msg := MockBroadcastMsg(i)
//				msgBytes, _ := json.Marshal(msg)
//
//				// 设置模拟的请求体
//				c.Request = httptest.NewRequest(http.MethodPost, "/broadcast", bytes.NewBuffer(msgBytes))
//				c.Request.Header.Set("Content-Type", "application/json")
//
//				// 调用GetIn函数
//				Broadcast(c)
//			}(i)
//		}
//
//		// 等待所有实例完成
//		wg.Wait()
//	}
func TestWs(t *testing.T) {
	// 定义一个在线用户列表
	// 定义一个在线用户列表
	var onlineUsers = []model.User{
		model.User{
			UserName:      "User1",
			Fingerprint:   "Fingerprint1",
			Host:          "Host1",
			LastHeartbeat: time.Now(),
		},
		model.User{
			UserName:      "User2",
			Fingerprint:   "Fingerprint2",
			Host:          "Host2",
			LastHeartbeat: time.Now(),
		},
		model.User{
			UserName:      "User3",
			Fingerprint:   "Fingerprint3",
			Host:          "Host3",
			LastHeartbeat: time.Now(),
		},
		model.User{
			UserName:      "User4",
			Fingerprint:   "Fingerprint4",
			Host:          "Host4",
			LastHeartbeat: time.Now(),
		},
		model.User{
			UserName:      "User5",
			Fingerprint:   "Fingerprint5",
			Host:          "Host5",
			LastHeartbeat: time.Now(),
		},
	}
	var requestMsg = communication.RequestMessage{
		Data:        communication.Data{},
		Content:     []communication.ContentReq{},
		OnlineUsers: communication.OnlineUsers{onlineUsers},
	}
	users, _ := json.Marshal(requestMsg)
	// 打印看看
	log.Println("OnlineUsers:", string(users))

}
