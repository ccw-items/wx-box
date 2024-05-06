package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

type dpHandlerType struct {
	Type              string `json:"type"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Time              string `json:"time"`
	NowPrice          string `json:"nowPrice"`
	YesterdayEndPrice string `json:"yesterdayEndPrice"`
	TodayStartPrice   string `json:"todayStartPrice"`
	TodayMax          string `json:"todayMax"`
	TodayMin          string `json:"todayMin"`
	Increase          string `json:"increase"`
	IncreasePerc      string `json:"increasePerc"`
	TraNumber         string `json:"traNumber"`
	TraAmount         string `json:"traAmount"`
	Chart             struct {
		Min   string `json:"min"`
		Day   string `json:"day"`
		Week  string `json:"week"`
		Month string `json:"month"`
	} `json:"chart"`
}

func dpHandler(msg *openwechat.Message, code ...string) {
	var stockCode string
	if len(code) > 0 {
		stockCode = code[0]
	} else {
		stockCode = "sh000001" // 默认股票代码
	}

	url := "暂时替换掉URL" + stockCode

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("大盘接口请求失败:", err)
		if msg != nil {
			msg.ReplyText("大盘接口请求失败")
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("大盘接口解析失败:", err)
		if msg != nil {
			msg.ReplyText("大盘接口解析失败")
		}
		return
	}

	var data dpHandlerType
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("大盘接口格式化失败:", err)
		if msg != nil {
			msg.ReplyText("大盘接口格式化失败")
		}
		return
	}

	if msg != nil {
		msg.ReplyText(fmt.Sprintf(
			"股票名称: %s(%s)\n当前价格: %s\n涨幅: %s\n涨幅百分比: %s%%\n更新时间: %s",
			data.Name,
			data.Code,
			data.NowPrice,
			data.Increase,
			data.IncreasePerc,
			data.Time,
		))
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
			v.SendText(
				fmt.Sprintf(
					"股票名称: %s(%s)\n当前价格: %s\n涨幅: %s\n涨幅百分比: %s%%\n更新时间: %s",
					data.Name,
					data.Code,
					data.NowPrice,
					data.Increase,
					data.IncreasePerc,
					data.Time,
				))
		}
	}
}
