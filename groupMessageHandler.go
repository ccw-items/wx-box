package main

import (
	"fmt"
	"time"

	"github.com/eatmoreapple/openwechat"

	"regexp"

	"strings"
)

func extractInvitee(str string) string {
	re := regexp.MustCompile(`邀请\"(.*)\"加入了群聊`)
	match := re.FindStringSubmatch(str)

	if len(match) > 0 {
		return match[1]
	} else {
		return ""
	}
}
func isMe(str string) (bool, string) {
	prefix := "@bots"
	if strings.HasPrefix(str, prefix) {
		return true, strings.TrimSpace(str[len(prefix):])
	} else {
		return false, ""
	}
}

// GroupMessageHandler 处理来自群组的消息
func GroupMessageHandler(msg *openwechat.Message) {
	if msg.IsSendByGroup() && msg.IsText() {
		groupSend, groupSendErr := msg.SenderInGroup()
		if groupSendErr != nil {
			fmt.Println("get groupSendErr :", groupSendErr)
			return
		}
		fmt.Println("群消息发送者是 :", groupSend)
		fmt.Println("群消息发送内容是 :", msg.Content)
		isPreFix, msgContent := isMe(msg.Content)
		strArray := []string{
			"----- [文字相关] -----",
			"*一言:获取1条一言",
			"*历史上的今天:获取历史上的今天",
			"*舔狗日记:获取一条舔狗日记",
			"*大盘:默认获取大A指数",
			"*sh:上海证券交易所代码（举例sh000001）",
			"*sz:深圳证券交易所代码（举例sz000001）",
			"*weibo | baidu | v2ex | qfc | zhihu | douban | toutiao | hupu | tianya | bilibili | maimai | ithome | zol | ifanr | oschina | csdn | huxiu | smzdm | douyin",
			"可获取相关网站热榜",
			"\n----- [图片相关] -----",
			"*随机一图:随机一张图片",
			"*self:打工人日历",
			"*大盘云图:获取大盘云图",
			"\n----- [数据处理] -----",
			"*gold:查询今日金价",
		}
		fmt.Println("判断命令是 :", msgContent)
		re := regexp.MustCompile(`^(sh|sz)\d{6}$`)
		if isPreFix {
			switch msgContent {
			case "菜单":
				msg.ReplyText(strings.Join(strArray, "\n"))
			case "一言":
				oneMsgHandler(msg)
			case "网抑云":
				fmt.Println("网抑云")
			case "gold":
				glodHandle(msg)
			case "随机一图":
				randomPictureHandle(msg)
			case "新闻":
				newsHander(msg)
			case "舔狗日记":
				tgrjHandler(msg)
			case "weibo", "baidu", "v2ex", "qfc", "zhihu", "douban", "toutiao", "hupu", "tianya", "bilibili", "maimai", "ithome", "zol", "ifanr", "oschina", "csdn", "huxiu", "smzdm", "douyin":
				hotHandler(msg, msgContent)
			case "self":
				selefMsgHandler(msg)
			case "大盘":
				dpHandler(msg)
			case "历史上的今天":
				ageToday(msg)
			case "大盘云图":
				dpDrawmHandler(msg)
			case "group":
				groupHandle()
			default:
				if re.MatchString(msgContent) {
					dpHandler(msg, msgContent)
				} else {
					gpthandle(msg, msgContent)
				}
			}
		} else {
			return
		}
	}

	if !msg.IsText() && extractInvitee(msg.Content) != "" {
		invitee := extractInvitee(msg.Content)
		fmt.Println("被邀请人是:", invitee)
		reply := fmt.Sprintf("恭喜 %s 加入群聊", invitee)
		time.Sleep(1 * time.Second) // 延迟5秒
		msg.ReplyText(reply)
	}
}
