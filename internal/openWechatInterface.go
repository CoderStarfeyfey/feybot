/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/21
 **/
package internal

import (
	"errors"
	"fmt"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	"github.com/eatmoreapple/openwechat"
	"os"
	"path"
	"regexp"
)

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
func SendText(msg *openwechat.Message, text string) error {
	_, err := msg.ReplyText(text)
	if err != nil {
		utils.FeyLog.Debugf("fail to send the text")
		return err
	}
	return nil

}
func SendPic(msg *openwechat.Message, path string) error {
	isExist := fileExist(path)
	if !isExist {
		utils.FeyLog.Debugf("file path %s do not exist", path)
		return errors.New("path is not exist")
	}
	content, err := os.Open(path)
	if err != nil {
		utils.FeyLog.Debugf("fail to load the resource")
		return errors.New("load jpg fail")
	}
	msg.ReplyImage(content)
	return nil
}
func SendTextAndPic(msg *openwechat.Message, result string) error {

	re := regexp.MustCompile(`path:(.*?) text:(.*)`)

	// 在字符串中查找匹配项
	matches := re.FindStringSubmatch(result)
	if len(matches) > 2 {
		// 提取出匹配的组
		picPath := matches[1]
		text := matches[2]
		picPath = path.Join(utils.PicDir, picPath)
		fmt.Println("Path:", picPath)
		fmt.Println("Text:", text)
		SendText(msg, text)
		SendPic(msg, picPath)
		return nil
	} else {
		fmt.Println("No match found")
	}
	return nil
}

func GetGroupName(msg *openwechat.Message) (string, string) {
	sender, err := msg.Sender()
	if err != nil {
		return "", "GetSender"
	}

	group, terror := sender.AsGroup() // 将sender转为group类型

	if terror != true || group == nil {
		return "", "Getgroup"
	}

	group.Detail() // 获取详细信息，保证信息及时更新
	return group.NickName, ""
}

func GetGroupObj(msg *openwechat.Message) (group *openwechat.Group, err error) {
	sender, err := msg.Sender()
	if err != nil {
		return nil, fmt.Errorf("Get sender fail")
	}

	group, terror := sender.AsGroup() // 将sender转为group类型
	if !terror {
		return nil, fmt.Errorf("tranfer sender to group fail")
	}
	return group, nil

}

func SendMsgToGroupByObj(group *openwechat.Group, content string) error {
	_, err := group.SendText(content)
	if err != nil {
		fmt.Errorf("Send msg to group fail ,Error :%v", err)
	}
	return nil
}
