package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

func newsHander(msg *openwechat.Message) {
	resp, err := http.Get("https://dayu.qqsuu.cn/weiyujianbao/apis.php")
	if err != nil {
		fmt.Print("获取每日新闻失败")
		if msg != nil {
			msg.ReplyText("获取每日新闻失败")
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("解析每日新闻失败")
		if msg != nil {
			msg.ReplyText("解析每日新闻失败")
		}
	}
	if msg != nil {
		reader := bytes.NewReader(body)
		msg.ReplyImage(reader)
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
			reader := bytes.NewReader(body)
			v.SendImage(reader)
		}
	}
}
