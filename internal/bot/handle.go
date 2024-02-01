/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/21
 **/
package bot

import (
	"fmt"
	"github.com/CoderStarfeyfey/feybot/config"
	"github.com/CoderStarfeyfey/feybot/internal"
	"github.com/CoderStarfeyfey/feybot/internal/plugins"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	"github.com/eatmoreapple/openwechat"
	"reflect"
	"strings"
)

var _ messageHandleInterface = (*GroupMessageHandler)(nil)

// 定义消息处理接口，无论是群聊还是私聊都需要实现该接口的方法
type messageHandleInterface interface {
	handle(msg *openwechat.Message) error
	reply(msg *openwechat.Message) error
}

// 判断消息是否命中默认的回复的消息列表
func isHandleDefaultReply(req string, normalReplyMap map[string]config.DefaultReplyStruct) (bool, string) {
	for key, value := range normalReplyMap {
		if utils.StringSubContains(value.Content, req) {
			return true, key
		}
	}
	return false, ""
}

// 通用请求处理过程
func CommonHandleRequest(req *internal.RequestStruct) (*internal.ReplyStruct, error) {
	//1.处理默认配置中已有的处理方式
	//2.处理插件中定义的方式
	//3.使用gpt3.5模型进行回复
	//这里用反射，运行时确定
	if res, msgType := isHandleDefaultReply(req.RequestTxt, config.NormalReplyMap); res == true {
		//如果是类型1就直接按照随机种子从列表中回复
		//如果是类型2直接回复字符串并且回复一张图片
		rv := reflect.ValueOf(config.FeyConfig.DefalutReplayMsg)
		if !rv.FieldByName(msgType).IsValid() {
			utils.FeyLog.Errorf("config  could not find the field %s", msgType)
			return nil, nil
		}
		fieldValue := rv.FieldByName(msgType)
		//反射取出对应的结构体
		value := fieldValue.Interface()
		//使用类型断言
		if coolMsg, ok := value.([]string); ok {
			index := utils.RandomLessThan(len(coolMsg))
			return &internal.ReplyStruct{1, coolMsg[index]}, nil
		} else if thankMsg, ok := value.(config.DefaultThankMsgStruct); ok {
			result := fmt.Sprintf("path:%s text:%s", thankMsg.ThankMsgPic, thankMsg.ThankMsgText)
			return &internal.ReplyStruct{3, result}, nil

		}
	}
	if value, ok := internal.DefaultReplyMap[req.RequestTxt]; ok {
		utils.FeyLog.Debugf("use the default answer to reply")
		return &internal.ReplyStruct{ReType: value.ReType, ReText: value.ReText}, nil
	}
	if feature, ok := config.FeyConfig.Features[req.RequestTxt]; ok {
		if !feature.Enable ||
			!feature.EnableGroupWbList[req.Groupname] {
			return &internal.ReplyStruct{ReType: 1, ReText: "此功能暂时无法使用,请联系管理员授权该功能或者开启功能"}, nil
		}
		//用户授权成功并且开启enable,调用插件
		if funcName, ok := internal.PluginMap[feature.EntryFunctionName]; ok {
			reply, err := funcName(req.Uname, req.Uuid)
			if err == nil {
				utils.FeyLog.Debugf("Find the plugin to use,and run it successfully")
				return reply, nil
			} else {
				utils.FeyLog.Errorf("fail to run the plugin,error:%s", err)
				return reply, err
			}
		}
	}
	//不是默认回复，也不是请求插件，全都用通用千问来回答
	reply, err := plugins.TYQWAutoReply(req)
	if err != nil {
		utils.FeyLog.Errorf("Request to tyqw fail,error:%v", err)
		return nil, err
	} else {
		return reply, nil
	}
	return nil, nil
}

// 定义群聊消息结构体
type GroupMessageHandler struct {
}

// 定义私聊消息结构体
type UserMessageHandler struct {
}

func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	groupName, err := internal.GetGroupName(msg)
	if err != "" {
		return nil
	}
	//判断群聊是否在白名单中，非白名单的不做回复
	if isExistWhiteList := config.FeyConfig.GroupWhiteListMap[groupName]; isExistWhiteList == false {
		utils.FeyLog.Debugf("%s do not have right to use the feybot", groupName)
		if msg.IsAt() {
			msg.ReplyText("该群并没有权限使用feybot功能哦，请联系feybot管理员开通")
		}
		return nil
	}
	//目前只开通了文本回复的功能，图片回复后续加入
	if msg.IsText() && msg.Content != "" {
		return g.reply(msg)
	}
	return nil
}
func (g *GroupMessageHandler) reply(msg *openwechat.Message) error {
	groupName, ferr := internal.GetGroupName(msg)
	if ferr != "" {
		return nil
	}
	utils.FeyLog.Debugf("receive the msg from %v,content:%v", groupName, msg.Content)
	//只处理@的消息
	if !msg.IsAt() {
		return nil
	}
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		utils.FeyLog.Errorf("error to get sender,error:%v", err)
		return nil
	}
	if groupSender.DisplayName == "" {
		groupSender.DisplayName = groupSender.NickName
	}
	// 替换掉@文本
	replaceText := "@" + config.FeyConfig.BotName
	requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	//临时插入之后由智能识别插件完成
	//if strings.Contains(replaceText, "帅哥") || strings.Contains(replaceText, "靓仔") {
	//
	//	return internal.SendText(msg)
	//}
	reply, err := CommonHandleRequest(&internal.RequestStruct{
		Uname:      groupSender.DisplayName,
		Uuid:       groupSender.UserName,
		Groupname:  groupName,
		RequestTxt: requestText,
	})
	if err != nil {
		//发送错误消息
		return nil
	}
	if reply.ReType == 1 {
		return internal.SendText(msg, reply.ReText)
	} else if reply.ReType == 2 {

		return internal.SendPic(msg, reply.ReText)
	} else if reply.ReType == 3 {
		return internal.SendTextAndPic(msg, reply.ReText)
	}
	return nil
}
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	return nil
}
func (g *UserMessageHandler) reply(msg *openwechat.Message) error {
	return nil
}

// 定义两个工厂函数，生成群聊和私聊的结构体对象给外部使用，所有的初始化过程均在这里实现
func NewGroupMessageHandle() messageHandleInterface {
	return &GroupMessageHandler{}
}
func NewUserMessageHandle() messageHandleInterface {
	return &UserMessageHandler{}
}

var HandlerMap = make(map[string]messageHandleInterface)

func init() {
	HandlerMap["groupHandler"] = NewGroupMessageHandle()
	HandlerMap["userHandler"] = NewUserMessageHandle()
}
func Handle(msg *openwechat.Message) {
	if msg.IsSendByGroup() {
		HandlerMap["groupHandler"].handle(msg)
		return
	} else if msg.IsSendByFriend() {
		HandlerMap["userHandler"].handle(msg)
		return
	}
}
