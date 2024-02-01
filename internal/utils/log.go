/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/20
 **/
package utils

import (
	"github.com/CoderStarfeyfey/feybot/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	_log "log"
	"os"
)

var FeyLog = logrus.New()
var logFile *lumberjack.Logger
var (
	logDefaultDir = "./botLog"
	TmpDir        = logDefaultDir + "/tmp" // 临时目录
	LogDir        = logDefaultDir + "/log" // 日志目录
	PicDir        = logDefaultDir + "/pic" // 图片目录
	LogFilename   = LogDir + "/feybot.log"
)

func initlogDir() {
	if _, err := os.Stat(LogDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDefaultDir, 0755)
		_log.Println("Successfully create the log dir %s", LogDir)
		if err != nil {
			_log.Fatalln("Create log dir {%s} failed", LogDir)
		}
	}
}
func FeyLogInit() {
	initlogDir()
	logFile = &lumberjack.Logger{
		Filename:   LogFilename,
		MaxSize:    config.FeyConfig.FeyLogConfig.Maxsize,    // 日志最大占用，单位MB
		MaxBackups: config.FeyConfig.FeyLogConfig.Maxbackups, //最多切片数量
		MaxAge:     config.FeyConfig.FeyLogConfig.Maxage,     //日志最多存活时间
		Compress:   config.FeyConfig.FeyLogConfig.Compress,   //是否压缩
	}

	// 日志默认级别设置为Info
	FeyLog.SetLevel(logrus.DebugLevel)
	// 记录日志时自动包含调用函数的信息
	FeyLog.SetReportCaller(true)
	// 日志输出到文件
	FeyLog.SetOutput(logFile)

	FeyLog.Infof("feybot init finished")
}
