package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

type TGRResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Form string `json:"form"`
}

func tgrjHandler(msg *openwechat.Message) {
	resp, err := http.Get("https://api.5qmn.com/api/tiangou/?type=json")
	if err != nil {
		fmt.Print("舔狗日记获取失败")
		if msg != nil {
			msg.ReplyText("舔狗日记获取失败")
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("舔狗日记解析失败")
		if msg != nil {
			msg.ReplyText("舔狗日记解析失败")
		}
	}
	var data TGRResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("舔狗日记格式化失败", err)
		if msg != nil {
			msg.ReplyText("舔狗日记格式化失败")
		}
	}

	if msg != nil {
		msg.ReplyText(fmt.Sprintf("%s \n ---《舔狗日记》", data.Msg))
	}
}
