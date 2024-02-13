package bot

import (
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	"github.com/eatmoreapple/openwechat"
)

func FeybotRun() {
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式
	// 注册消息处理函数
	bot.MessageHandler = Handle

	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// 执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		if err = bot.Login(); err != nil {
			utils.FeyLog.Errorf("login error: %v", err)
			return
		}
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
