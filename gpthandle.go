package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eatmoreapple/openwechat"
)

const (
	apiKey  = "sk-xxxxxxxxxxxx"                                   //我的
	apiUrl  = "https://api.chatanywhere.tech/v1/chat/completions" // 修改为新的请求地址
	modelId = "gpt-3.5-turbo"                                     // ChatGPT模型ID
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model        string    `json:"model"`
	Messages     []Message `json:"messages"`
	Instructions string    `json:"instructions"`
	Description  string    `json:"description"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Logprobs     string  `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

func gpthandle(msg *openwechat.Message, msgContent string) {
	requestData := Request{
		Model: modelId,
		Messages: []Message{
			{Role: "user", Content: msgContent},
		},
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		msg.ReplyText(fmt.Sprintf("GPT请求数据编码失败: %s", err))
		return
	}

	fmt.Println("Request Body:", string(jsonData))

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		msg.ReplyText(fmt.Sprintf("创建请求失败: %s", err))
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		msg.ReplyText(fmt.Sprintf("GPT请求失败: %s", err))
		return
	}
	defer resp.Body.Close()

	rawResponseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		msg.ReplyText(fmt.Sprintf("读取响应体失败: %s", err))
		return
	}
	fmt.Println("Raw Response Body:", string(rawResponseBody))

	// 解析响应体以检查是否有错误
	var responseErr struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	err = json.Unmarshal(rawResponseBody, &responseErr)
	if err == nil && responseErr.Error.Type == "chatanywhere_error" && responseErr.Error.Code == "429 TOO_MANY_REQUESTS" {
		msg.ReplyText("每日API达到上限")
		return
	}

	// 如果没有错误，继续处理响应
	var responseData Response
	err = json.Unmarshal(rawResponseBody, &responseData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		msg.ReplyText(fmt.Sprintf("解析响应数据失败: %s", err))
		return
	}

	if len(responseData.Choices) > 0 {
		groupSend, groupSendErr := msg.SenderInGroup()
		if groupSendErr != nil {
			fmt.Println("get groupSendErr :", groupSendErr)
			return
		}
		if responseData.Choices[0].Message.Content != "" {
			msg.ReplyText(fmt.Sprintf("@%s\n %s", groupSend.NickName, responseData.Choices[0].Message.Content))
		} else {
			msg.ReplyText(fmt.Sprintf("@%s\n %s", groupSend.NickName, "触发了GPT的过滤，换个话题"))
		}
	} else {
		fmt.Println("没有收到选择")
	}
}
