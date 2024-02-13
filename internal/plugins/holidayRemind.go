/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/21
 **/
package plugins

import (
	"github.com/CoderStarfeyfey/feybot/internal"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	"os"
)

var CalenderUrl = "https://api.vvhan.com/api/moyu?type=json"

type Cal_data struct {
	Success bool   `json:"success`
	Url     string `json:"url`
}

func init() {
	internal.PluginMap["HolidaysQuery"] = HolidaysQuery
}
func HolidaysQuery(req *internal.RequestStruct) (*internal.ReplyStruct, error) {
	tmpFilePath, _ := utils.GenPicTempPath()
	if _, err := os.Stat(tmpFilePath); err == nil {
		utils.FeyLog.Debugf("Path : %v has been there,use cache to upload")
		return &internal.ReplyStruct{
			ReType: 2,
			ReText: tmpFilePath,
		}, nil
	}
	var calData Cal_data
	err := utils.GetHttpRespJson(CalenderUrl, &calData)
	if err != nil {
		utils.FeyLog.Debugf("Send req to fish cal fail ,error:%v", err)
		return nil, err
	}
	err = utils.DownloadPicToDefaultDir(calData.Url, tmpFilePath)
	if err != nil {
		utils.FeyLog.Debugf("Down load holiday remind fail,error:%v", err)
		return nil, err
	}
	return &internal.ReplyStruct{
		ReType: 2,
		ReText: tmpFilePath,
	}, nil
}
