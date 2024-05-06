package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
)

// MessageHandler 注册消息处理函数
func messageHander(bot *openwechat.Bot) {
	bot.MessageHandler = func(msg *openwechat.Message) {
		FriendMessageHandler(msg)
		GroupMessageHandler(msg)

		if msg.IsText() && msg.Content == "ping" {
			fmt.Println(msg.Content)
			msg.ReplyText("pong")
		}
	}
}
