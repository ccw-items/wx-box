package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

type mgsType struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func hotHandler(msg *openwechat.Message, msgType string) {
	hotListMap := map[string]string{
		"weibo":    "微博前10",
		"baidu":    "百度前10",
		"v2ex":     "V2EX前10",
		"qfc":      "晴风村前10",
		"zhihu":    "知乎前10",
		"douban":   "豆瓣前10",
		"toutiao":  "今日头条前10",
		"hupu":     "虎扑前10",
		"tianya":   "天涯前10",
		"bilibili": "哔哩哔哩前10",
		"maimai":   "脉脉前10",
		"ithome":   "IT之家前10",
		"zol":      "中关村在线前10",
		"ifanr":    "爱范儿前10",
		"oschina":  "开源中国前10",
		"csdn":     "CSDN前10",
		"huxiu":    "虎嗅网前10",
		"smzdm":    "什么值得买前10",
		"douyin":   "抖音前10",
	}
	hotListName, ok := hotListMap[msgType]
	if !ok {
		msg.ReplyText("未找到相应的热榜")
		return
	}
	var url = `https://api.lancely.tech/hot/` + msgType
	resp, err := http.Get(url)
	if err != nil {
		msg.ReplyText("信息接口请求失败")
		fmt.Println("报错", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		msg.ReplyText("信息接口解析失败")
		return
	}
	var data []mgsType
	err = json.Unmarshal(body, &data)
	if err != nil {
		msg.ReplyText("请求格式化失败")
		return
	}
	var results []string
	count := 0
	for _, item := range data {
		if count >= 10 {
			break
		}
		result := fmt.Sprintf("%d:%s\n%s", count+1, item.Title, item.Link)
		results = append(results, result)
		count++
	}
	additionalText := fmt.Sprintf("🔥🔥🔥%s的热榜\n", hotListName)
	replyText := additionalText + strings.Join(results, "\n")
	msg.ReplyText(replyText)
}
