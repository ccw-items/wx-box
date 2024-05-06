package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

type Event struct {
	Year  string `json:"year"`
	Title string `json:"title"`
}

type ageTodayType struct {
	Today  string  `json:"today"`
	Result []Event `json:"result"`
}

func ageToday(msg *openwechat.Message) {
	resp, err := http.Get("https://www.ipip5.com/today/api.php?type=json")
	if err != nil {
		fmt.Print("历史上的今天获取失败")
		if msg != nil {
			msg.ReplyText("历史上的今天获取失败")
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("解析历史上的今天失败")
		if msg != nil {
			msg.ReplyText("解析历史上的今天失败")
		}
	}
	var events ageTodayType
	err = json.Unmarshal(body, &events)

	if err != nil {
		fmt.Println("格式化历史上的今天失败", err)
		if msg != nil {
			msg.ReplyText("格式化历史上的今天失败")
		}
	}

	var eventStrings []string
	eventStrings = append(eventStrings, "时间:"+events.Today)
	for _, event := range events.Result {
		if event.Year != "2024" {
			eventStrings = append(eventStrings, event.Year+":"+event.Title)
		}
	}

	if msg != nil {
		msg.ReplyText(strings.Join(eventStrings, "\n"))
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
			v.SendText(strings.Join(eventStrings, "\n"))
		}
	}
}
