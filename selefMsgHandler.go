package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

// selefMsgHandler 注册消息处理函数
func selefMsgHandler(msg *openwechat.Message) {

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

	resp, err := http.Get("https://dayu.qqsuu.cn/moyuribao/apis.php")
	if err != nil {
		fmt.Println("摸鱼日历接口请求失败")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("摸鱼日历接口解析失败")
		return
	}
	if msg == nil {
		for _, v := range grp {
			reader := bytes.NewReader(body)
			v.SendImage(reader)
		}
	} else {
		reader := bytes.NewReader(body)
		msg.ReplyImage(reader)
	}
}
