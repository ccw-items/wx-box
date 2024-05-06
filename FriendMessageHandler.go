package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
)

// FriendMessageHandler 处理来自朋友的消息
func FriendMessageHandler(msg *openwechat.Message) {
	if msg.IsSendByFriend() && msg.IsText() {
		send, sendErr := msg.Sender()
		if sendErr != nil {
			fmt.Println("get sendError :", sendErr)
			return
		}
		fmt.Println("个人消息发送者 发送人是:", send)
		fmt.Println("个人消息发送者 发送信息是是:", msg.Content)
	}
}
