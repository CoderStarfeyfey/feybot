/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/23
 **/
package plugins

import (
	"context"
	"fmt"
	"github.com/CoderStarfeyfey/feybot/config"
	"github.com/CoderStarfeyfey/feybot/internal"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	openai "github.com/sashabaranov/go-openai"
)

const maxQueryGPTTimes = 10
const chatHistoryCache = 5

type ChatQuesAndReply struct {
	Question string
	Answer   string
}

var client *openai.Client

// 记录对应用户已经使用gpt的次数，一个账号只能使用5次。升级成feybot plus++可以一天使用15次
var UserUseGptTimesMap = make(map[string]int)

// 记录该用户使用回复记录缓存
var UserUseGptHistory = make(map[string][]ChatQuesAndReply)

// 记录该消息对应的时间，5分钟清理缓存一次
var UserUseGptTime = make(map[string]int64)

// 使用泛型实现的 moveElementToEnd
func moveElementToEnd[T any](slice []T, index int) []T {
	// 检查索引的有效性
	if index < 0 || index >= len(slice) {
		return slice
	}

	// 删除索引处的元素并保存这个元素
	element := slice[index]
	slice = append(slice[:index], slice[index+1:]...)

	// 将元素追加到切片末尾
	slice = append(slice, element)

	return slice
}

func init() {
	client = openai.NewClient(config.FeyConfig.ChatGPTToken)
}
func GPTAutoReply(req *internal.RequestStruct) (*internal.ReplyStruct, error) {
	//这里要实现一个缓存机制以及除了管理员用户以外最大用gpt的访问次数
	//先判断已经使用gpt的次数
	if UserUseGptTimesMap[req.Uuid] == maxQueryGPTTimes {
		utils.FeyLog.Debugf("User:%v has require the maxmium request today")
		return &internal.ReplyStruct{
			ReType: 1,
			ReText: "you has run out of the times to query gpt,please use plus+ feybot",
		}, nil
	}
	//在判断是否有缓存，如果有缓存返回该值并且刷新该历史记录的时间，重新排序这个列表,
	if value, ok := UserUseGptHistory[req.Uuid]; ok {
		for i, v := range value {
			if v.Question == req.RequestTxt {
				utils.FeyLog.Debugf("Find the history "+
					"cache user:%v,"+
					"question:%v"+
					"reply:%v", req.Uname, v.Question, v.Answer)
				UserUseGptHistory[req.Uuid] = moveElementToEnd(value, i)
				return &internal.ReplyStruct{
					ReType: 1,
					ReText: v.Answer,
				}, nil
			}
		}
	}
	//这里只实现了无状态回复，没有将历史记录的消息一起回复。
	if reply, err := CreateChatCompletion(req.RequestTxt); err != nil {
		utils.FeyLog.Errorf("chatGPT get reply error,error:%v", err)
		return nil, err
	} else {
		//用户新增次数+1
		UserUseGptTimesMap[req.Uuid] = UserUseGptTimesMap[req.Uuid] + 1
		//把这一次的问答加入到缓存中，这里需要判断是否缓存已经满了
		if len(UserUseGptHistory[req.Uuid]) == chatHistoryCache {
			utils.FeyLog.Debugf("user :%v clear the history", req.Uname)
			UserUseGptHistory[req.Uuid] = UserUseGptHistory[req.Uuid][1:]
		}
		UserUseGptHistory[req.Uuid] = append(UserUseGptHistory[req.RequestTxt], ChatQuesAndReply{
			Question: req.RequestTxt,
			Answer:   reply,
		})
		return &internal.ReplyStruct{
			ReType: 1,
			ReText: reply,
		}, nil
	}
}
func CreateChatCompletion(message string) (string, error) {
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)
	if err != nil {
		fmt.Println("发生错误:", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
