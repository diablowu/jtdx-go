package qywx

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendAgentMessage(agentID int, userID string, message string) {

	url := QYAPIEndpoint + "/cgi-bin/message/send?access_token=" + accessToken.Token

	tm := TextMessage{
		To:      userID,
		Type:    "text",
		AgentID: agentID,
		Text: struct {
			Content string `json:"content"`
		}{
			Content: message,
		},
	}

	if bs, err := json.Marshal(tm); err == nil {
		http.Post(url, "application/json", bytes.NewBuffer(bs))
	}

}

type TextMessage struct {
	To      string `json:"touser"`
	Type    string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}
