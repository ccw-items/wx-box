package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
	"github.com/robfig/cron/v3"
)

var Bot *openwechat.Bot

func main() {

	c := cron.New()

	Bot = openwechat.DefaultBot(openwechat.Desktop)
	messageHander(Bot)
	Bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	if err := Bot.Login(); err != nil {
		fmt.Println(err)
		return
	}

	self, err := Bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	c.AddFunc("00 9 * * *", func() {
		selefMsgHandler(nil)
	})

	c.AddFunc("00 10 * * *", func() {
		newsHander(nil)
	})

	c.AddFunc("00 23 * * *", func() {
		oneMsgHandler(nil)
	})

	c.AddFunc("0/1 * * * *", func() {
		groupHandle()
	})

	groups, err := self.Groups()

	c.AddFunc("30 18 * * 1-5", func() {
		fmt.Println("执行了这里")
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
			v.SendText("请赶紧下班不要浪费公司电。")
		}
	})

	// c.AddFunc("35 9 * * 1-5", func() {
	// 	dpHandler(nil)
	// })
	// c.AddFunc("35 11 * * 1-5", func() {
	// 	dpHandler(nil)
	// })
	// c.AddFunc("05 13 * * 1-5", func() {
	// 	dpHandler(nil)
	// })
	// c.AddFunc("05 15 * * 1-5", func() {
	// 	dpHandler(nil)
	// })

	c.Start()

	fmt.Println(groups, err)
	Bot.Block()
}
