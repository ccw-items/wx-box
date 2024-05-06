package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eatmoreapple/openwechat"
)

type GoldPrice struct {
	BrandName     string `json:"brandName"`
	ProductName   string `json:"productName"`
	Price         string `json:"price"`
	PriceUnit     string `json:"priceUnit"`
	RaiseDownType string `json:"raiseDownType"`
	UpdateDate    string `json:"updateDate"`
}

func glodHandle(msg *openwechat.Message) {
	resp, err := http.Get("暂时替换掉URL")
	if err != nil {
		msg.ReplyText("金价接口请求失败")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		msg.ReplyText("金价接口解析失败")
		return
	}

	var data []GoldPrice
	err = json.Unmarshal(body, &data)
	if err != nil {
		msg.ReplyText("请求格式化失败")
		return
	}

	var results []string
	for _, price := range data {
		result := fmt.Sprintf("%s:%s%s", price.BrandName, price.Price, price.PriceUnit)
		results = append(results, result)
	}

	// 检查data是否为空，避免出现下标越界的错误
	if len(data) > 0 {
		// 在最后添加时间信息
		results = append(results, "时间:"+data[0].UpdateDate)
	}

	msg.ReplyText(strings.Join(results, "\n"))
}
