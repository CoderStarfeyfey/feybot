/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/20
 **/
package config

import (
	"encoding/json"
	_log "log"
	"os"
)

type featureStruct struct { //定义插件的结构体
	Enable            bool            `json:"enable"`
	EntryFunctionName string          `json:"entryFunctionName"`
	EnableGroupWbList map[string]bool `json:"enableGroupWbList"`
}
type DefalutReplyMsgStruct struct {
	CoolManMsg []string              `json:"coolManMsg"` //夸奖类回复
	ErrorMsg   string                `json:"errorMsg"`
	ThankMsg   DefaultThankMsgStruct `json:"thankMsg"`
}
type DefaultThankMsgStruct struct {
	ThankMsgText string `json:"thankMsgText"`
	ThankMsgPic  string `json:"thankMsgPic"`
}
type FeyConfigStruct struct {
	BotName          string                `json:"botName"` //机器人名称
	DataDir          string                `json:"dataDir"` //机器人数据存放目录
	TyqwToken        string                `json:"tyqwToken"`
	ChatGPTToken     string                `json:"chatGPTToken"`
	DefalutReplayMsg DefalutReplyMsgStruct `json:"defalutReplayMsg"`
	GroupWhiteList   []string              `json:"groupWhiteList"` //群白名单，只有白名单的群才有特定的功能
	FeyLogConfig     struct {              //日志输出文件的配置
		Maxsize    int  `json:"maxsize"`
		Maxbackups int  `json:"maxbackups"`
		Maxage     int  `json:"maxage"`
		Compress   bool `json:"compress"`
	} `json:"feyLogConfig"`
	Features          map[string]featureStruct `json:"feature"`           //定义了机器人所使用的插件
	GroupWhiteListMap map[string]bool          `json:"groupWhiteListMap"` //群聊白名单Map
}
type DefaultReplyStruct struct {
	Content  []string `json:"content"`
	ReLength int      `json:"reLength"`
	ReType   int      `json:"reType"`
}
type NormalReplyStruct struct {
	CoolManMsg DefaultReplyStruct `json:"coolManMsg"`
	ThankMsg   DefaultReplyStruct `json:"thankMsg"`
}

var NormalReplyMap = make(map[string]DefaultReplyStruct)

// 这里定义一个全局结构体指针
var FeyConfig *FeyConfigStruct = new(FeyConfigStruct)
var NormalReplyConfig *NormalReplyStruct = new(NormalReplyStruct)

const ConfigDefaultPath = "./config/feybot.config"
const ConfigNormalReplyPath = "./config/normalReplyConfig.config"

func LoadConfig(paths ...string) {
	var configPath string
	if len(paths) == 0 {
		configPath = ConfigDefaultPath
	} else {
		configPath = paths[0]
	}
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		_log.Fatalln("map the config file to struct failed ", err)
	}
	if err = json.Unmarshal(configContent, FeyConfig); err != nil {
		_log.Fatalln("map the config file to struct failed ", err)
	}
	FeyConfig.GroupWhiteListMap = map[string]bool{}
	for _, groupName := range FeyConfig.GroupWhiteList {
		FeyConfig.GroupWhiteListMap[groupName] = true
	}
	_log.Println("Successfully parsed successfully the feybot file")
	replyContent, err := os.ReadFile(ConfigNormalReplyPath)
	if err != nil {
		_log.Fatalln("map the config file to struct failed ", err)
	}
	if err = json.Unmarshal(replyContent, &NormalReplyMap); err != nil {
		_log.Fatalln("map the config file to struct failed ", err)
	}
	_log.Println("Successfully parsed successfully the normalReplyConfig file")
}
