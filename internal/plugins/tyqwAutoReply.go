/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/24
 **/
package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CoderStarfeyfey/feybot/config"
	"github.com/CoderStarfeyfey/feybot/internal"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var TyqwApiUrl = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"

type TyqwResponse struct {
	Output struct {
		FinishReason string `json:"finish_reason"`
		Text         string `json:"text"`
	} `json:"output"`
}
type TyqwMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TyqwInput struct {
	Model string `json:"model"`
	Input struct {
		Prompt string `json:"prompt"`
	} `json:"input"`
}
type TyqwHistoryStruct struct {
	Question string
	Answer   string
}

func (ths TyqwHistoryStruct) String() string {
	return fmt.Sprintf("Question: %s, Answer: %d", ths.Question, ths.Answer)
}

const maxHistoryCacheLength = 10
const maxQueryTimes = 3

var UserHistoryCache = make(map[string][]TyqwHistoryStruct)
var UserQueryTimes = make(map[string]int)

func TYQWAutoReply(req *internal.RequestStruct) (*internal.ReplyStruct, error) {
	if _, ok := UserQueryTimes[req.Uuid]; !ok {
		UserQueryTimes[req.Uuid] = 0
	}
	if UserQueryTimes[req.Uuid] == maxQueryTimes {
		utils.FeyLog.Debugf("user:%v are unable to use AI model to reply", req.Uname)
		return &internal.ReplyStruct{
			ReType: 1,
			ReText: "你能够使用AI智能问答的功能已经超限,请升级feybot plus++",
		}, nil
	}
	//先判断缓存里是否存在，如果存在取出排序，如果不存在就添加到缓存区中
	if value, ok := UserHistoryCache[req.Uuid]; ok {
		for index, last := range value {
			if last.Question == req.RequestTxt {
				lastAnswer := last.Answer
				moveElementToEnd(value, index)
				return &internal.ReplyStruct{
					ReType: 1,
					ReText: lastAnswer,
				}, nil
			}
		}
		//用户存在，但是没有该历史缓存
		if len(UserHistoryCache[req.Uuid]) == maxHistoryCacheLength {
			UserUseGptHistory[req.Uuid] = UserUseGptHistory[req.Uuid][1:]
		}
	} else {
		//该用户第一次访问建立一个空切片
		UserHistoryCache[req.Uuid] = []TyqwHistoryStruct{}
	}
	reply, err := getTYQWReply(req.RequestTxt)
	if err != nil {
		return nil, err
	}
	UserHistoryCache[req.Uuid] = append(UserHistoryCache[req.Uuid], TyqwHistoryStruct{
		Question: req.RequestTxt,
		Answer:   reply,
	})
	UserQueryTimes[req.Uuid] += 1
	return &internal.ReplyStruct{
		ReType: 1,
		ReText: reply,
	}, nil
}
func getTYQWReply(content string) (string, error) {
	data := TyqwInput{
		Model: "qwen-turbo",
		Input: struct {
			Prompt string `json:"prompt"`
		}{
			Prompt: content,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// 创建一个请求
	req, err := http.NewRequest("POST", TyqwApiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+config.FeyConfig.TyqwToken)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 如果是token错误，视为没有配置，返回默认回复
	if resp.StatusCode == 401 {
		return "", errors.New("invalid token")
	}

	// 如果响应码不是200，打日志并返回error
	if resp.StatusCode != 200 {
		return "", errors.New("请求成功，但响应失败")
	}

	var response TyqwResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}

	return response.Output.Text, nil
}
func QueryUserTimes(uuid string) (string, error) {
	if times, ok := UserQueryTimes[uuid]; ok {
		return strconv.Itoa(times), nil
	}
	return "error to find times", fmt.Errorf("没有找到用户 %s 的次数", uuid)
}
func QueryUserHistory(uuid string) (string, error) {
	if res, ok := UserHistoryCache[uuid]; ok {
		var sb strings.Builder
		for _, entry := range res {
			sb.WriteString(entry.String())
			sb.WriteString(" | ") // 在每个元素之间添加分隔符
		}
		resultString := sb.String()
		if len(resultString) > 0 {
			resultString = resultString[:len(resultString)-3]
		}
		return resultString, nil
	}
	return "error to find history", fmt.Errorf("没有找到用户 %s 的历史缓存", uuid)
}

func GetUsersUUid(user string) (string, error) {
	if user != "user" {
		return "invalid args", fmt.Errorf("invliad args")
	}
	//考虑到字符串拼接方式性能的影响，这里用StringBuilder不用+=
	var builder strings.Builder
	for uuid, _ := range UserQueryTimes {
		builder.WriteString(uuid)
		builder.WriteString("\n")
	}
	res := builder.String()
	if len(res) == 0 {
		return "Do not have any uuids in cache", nil
	}
	return res, nil

}

var commandTyqwMap = make(map[string]func(uuid string) (string, error))

// 提供查询用户使用AI模型的次数
func RegistTyqwCommand(args []string) (string, error) {
	//feybot tyqw -u uuid -h 打印该用户的历史缓存
	//feybot tyqw -u uuid -t 打印该已经使用的次数
	//feybot tyqw -u user -q
	if len(args) <= 2 {
		return "", errors.New("args invalid")
	}
	funcName := commandTyqwMap[args[2]]
	res, err := funcName(args[1])
	if err != nil {
		return "", err
	}
	return res, nil
}

func init() {
	commandTyqwMap["-t"] = QueryUserTimes
	commandTyqwMap["-h"] = QueryUserHistory
	commandTyqwMap["-q"] = GetUsersUUid
	internal.RegisterCommand("tyqw", internal.CommandHandler{
		Execute: RegistTyqwCommand,
		Help: "tyqw [-t] - print out the times\n" +
			"tyqw [-u] - input uuid of user\n" +
			"tyqw [-h] - print history cache of user\n" +
			"tyqw [-q] - query uuids of all users\n",
	})
}
