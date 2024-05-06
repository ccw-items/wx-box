package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

// 图片数据结构
type ImageData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Total int `json:"total"`
		List  []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
			URL   string `json:"url"`
			Type  string `json:"type"`
		} `json:"list"`
	} `json:"result"`
}

// 处理随机图片请求
func randomPictureHandle(msg *openwechat.Message) {
	// 发送 HTTP 请求获取图片数据
	resp, err := http.Get("https://api.apiopen.top/api/getImages?page=0&size=1")
	if err != nil {
		msg.ReplyText("二次元图片获取失败")
		return
	}
	defer resp.Body.Close()

	// 解析 JSON 响应
	var imageData ImageData
	err = json.NewDecoder(resp.Body).Decode(&imageData)
	if err != nil {
		msg.ReplyText("二次元图片解析失败")
		return
	}

	// 检查是否成功获取到图片数据
	if imageData.Code != 200 {
		msg.ReplyText("二次元图片获取失败")
		return
	}

	// 获取图片 URL
	if len(imageData.Result.List) == 0 {
		msg.ReplyText("未找到图片")
		return
	}
	imageURL := imageData.Result.List[0].URL

	// 发送图片消息
	resp, err = http.Get(imageURL)
	if err != nil {
		msg.ReplyText("获取图片失败")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("解析随机一图失败")
		if msg != nil {
			msg.ReplyText("解析随机一图失败")
		}
	}
	reader := bytes.NewReader(body)
	msg.ReplyImage(reader)
}
