/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/31
 **/
package plugins

import (
	"errors"
	"fmt"
	"github.com/CoderStarfeyfey/feybot/internal"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
)

var LoveUrl = "https://api.vvhan.com/api/love?type=json"

type LoveWordsStruct struct {
	Success bool   `json:"success`
	Ishan   string `json:"ishan`
}

func init() {
	internal.PluginMap["LoveWords"] = LoveWords
}

func Test2() {
	var loveResp LoveWordsStruct
	err := utils.GetHttpRespJson(LoveUrl, &loveResp)
	if err != nil {
		utils.FeyLog.Debugf("Send req to love words fail ,error:%v", err)
	}
	fmt.Sprintf(loveResp.Ishan)
}
func LoveWords(uname, uuids string) (*internal.ReplyStruct, error) {
	var loveResp LoveWordsStruct
	err := utils.GetHttpRespJson(LoveUrl, &loveResp)
	if err != nil {
		utils.FeyLog.Debugf("Send req to love words fail ,error:%v", err)
		return nil, err
	}
	if loveResp.Success == false {
		utils.FeyLog.Debugf("Send req to love words fail,error:%v", loveResp.Ishan)
		return nil, errors.New("success :fail")
	}
	if err != nil {
		utils.FeyLog.Debugf("Down load holiday remind fail,error:%v", err)
		return nil, err
	}
	return &internal.ReplyStruct{
		ReType: 1,
		ReText: loveResp.Ishan,
	}, nil
}
