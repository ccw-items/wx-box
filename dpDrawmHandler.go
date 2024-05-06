package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

func dpDrawmHandler(msg *openwechat.Message) {
	resp, err := http.Get("暂时替换掉URL")
	if err != nil {
		msg.ReplyText("大盘云图获取失败")
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		msg.ReplyText("大盘云图解析失败")
		return
	}

	reader := bytes.NewReader(body)
	msg.ReplyImage(reader)
}
