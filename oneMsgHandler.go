package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

type ResponseData struct {
	Status string `json:"status"`
	Data   struct {
		Content string `json:"content"`
		Origin  string `json:"origin"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func oneMsgHandler(msg *openwechat.Message) {
	resp, err := http.Get("https://paul.ren/api/say")
	if err != nil {
		if msg != nil {
			msg.ReplyText("一言请求接口失败")
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if msg != nil {
			msg.ReplyText("一言解析接口失败")
		}
		return
	}

	var data ResponseData
	err = json.Unmarshal(body, &data)
	if err != nil {
		if msg != nil {
			msg.ReplyText("一言格式化接口失败")
		}
		return
	}
	result := fmt.Sprintf("%s - 《%s》", data.Data.Content, data.Data.Origin)
	if msg != nil {
		msg.ReplyText(result)
	} else {
		self, err := Bot.GetCurrentUser()
		if err != nil {
			fmt.Println("获取个人信息失败:", err)
			return
		}
		grp, err := self.Groups()
		if err != nil {
			fmt.Println("获取群组型系失败:", err)
			return
		}
		for _, v := range grp {
			v.SendText(result)
		}
	}
}
