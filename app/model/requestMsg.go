package model

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

type RequestMessage struct {
	Data    `json:"data"`
	Content []ContentReq `json:"content"`
}
