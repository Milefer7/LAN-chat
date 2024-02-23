package communication

// 定义通信返回数据结构

type ContentRes struct {
	Status  int8   `json:"status"`
	Message string `json:"message"`
}

// ResponseMessage 响应消息结构
type ResponseMessage struct {
	Data       `json:"data"`
	ContentRes `json:"content"`
}
