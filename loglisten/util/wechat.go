/*
createTime: 2022/12/3
*/
package util

import (
	"encoding/json"
	"net/http"
	"strings"
)

type wechatBotMsg struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content             string   `json:"content"`
		MentionedMobileList []string `json:"mentioned_mobile_list"`
	} `json:"text"`
}

func SendMsg(msg string, toUser []string) int {
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=a989637f-844c-4129-a610-0f1583ebc87f"
	method := "POST"

	msgInfo := wechatBotMsg{
		MsgType: "text",
		Text: struct {
			Content             string   `json:"content"`
			MentionedMobileList []string `json:"mentioned_mobile_list"`
		}{
			Content:             msg,
			MentionedMobileList: toUser,
		},
	}
	msgInfoByte, _ := json.Marshal(msgInfo)
	payload := strings.NewReader(string(msgInfoByte))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return 400
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return 401
	}
	return res.StatusCode
}
