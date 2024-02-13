package main

import (
	"github.com/CoderStarfeyfey/feybot/config"
	"github.com/CoderStarfeyfey/feybot/internal"
	"github.com/CoderStarfeyfey/feybot/internal/bot"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
)

func main() {
	//加载配置文件
	config.LoadConfig()
	//日志模块初始化
	utils.FeyLogInit()
	//开启日志的热调试功能
	go utils.HotDebugInit()
	//开启监听控制台调试功能
	go internal.DbgConsoleServiceHandle()
	//初始化机器人模块
	bot.FeybotRun()
}
