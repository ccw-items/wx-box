package main

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
)

type GroupMembers struct {
	GroupName      string   `json:"group_name"`
	PeopleNameList []string `json:"people_name_list"`
}

var groupMembersList []GroupMembers

func groupHandle() {
	self, err := Bot.GetCurrentUser()
	if err != nil {
		fmt.Println("获取个人信息失败:", err)
		return
	}
	groups, err := self.Groups()
	if err != nil {
		fmt.Println("获取群组信息失败:", err)
		return
	}
	for _, v := range groups {
		members, err := v.Members()
		if err != nil {
			fmt.Println("获取群成员列表失败:", err)
			return
		}
		isNewGroup := true
		for i, groupMember := range groupMembersList {
			if groupMember.GroupName == v.NickName {
				isNewGroup = false
				if len(members) != len(groupMember.PeopleNameList) {
					missingMembers := findMissingMembers(groupMember.PeopleNameList, members)
					fmt.Printf("群组 %s 中少了成员: %v\n", groupMember.GroupName, missingMembers)
					for _, missingMember := range missingMembers {
						v.SendText(fmt.Sprintf("检测到群友 %s 退出群聊", missingMember))
					}
					groupMembersList[i].PeopleNameList = []string{}
					for _, member := range members {
						groupMembersList[i].PeopleNameList = append(groupMembersList[i].PeopleNameList, member.NickName)
					}
				}
				break
			}
		}
		if isNewGroup {
			groupMembers := GroupMembers{
				GroupName:      v.NickName,
				PeopleNameList: make([]string, len(members)),
			}
			for i, member := range members {
				groupMembers.PeopleNameList[i] = member.NickName
			}
			groupMembersList = append(groupMembersList, groupMembers)
		}
	}
}

func findMissingMembers(oldMembers []string, newMembers openwechat.Members) []string {
	missingMembers := []string{}
	for _, oldMember := range oldMembers {
		found := false
		for _, newMember := range newMembers {
			if oldMember == newMember.NickName {
				found = true
				break
			}
		}
		if !found {
			missingMembers = append(missingMembers, oldMember)
		}
	}
	return missingMembers
}
