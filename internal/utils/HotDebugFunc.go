/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/26
 **/
package utils

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"os"
)

var debugFlagPath = TmpDir + "/Debug"

// 这个脚本主要实现日志热调试的功能
func checkDebugFlag() bool {
	if _, err := os.Stat(debugFlagPath); err != nil {
		return false
	}
	return true
}

func HotDebugInit() error {
	flagIsExist := checkDebugFlag()
	if flagIsExist {
		FeyLog.Infof("Debug flag exist,need to set log level to debug")
		FeyLog.SetLevel(logrus.DebugLevel)
	}
	if _, err := os.Stat(TmpDir); os.IsNotExist(err) {
		_, err = os.Create(TmpDir)
		if err != nil {
			FeyLog.Infof("Error to create dir:%s ,error:%v", TmpDir, err)
		}
	}
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		FeyLog.Errorf("Create new fsnotify failed, err info: %v", err)
		return err
	}
	defer watch.Close()

	if err := watch.Add(TmpDir); err != nil {
		FeyLog.Errorf("Watch debug flag add failed, err info: %v", err)
		return err
	}
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create && (checkDebugFlag() == true) {
						FeyLog.Info("Debug flag created, change log level to DebugLevel")
						FeyLog.SetLevel(logrus.DebugLevel)
					}
					if (ev.Op&fsnotify.Remove == fsnotify.Remove || ev.Op&fsnotify.Rename == fsnotify.Rename) && checkDebugFlag() == false {
						FeyLog.Info("Debug flag lossed, change log level to InfoLevel")
						FeyLog.SetLevel(logrus.InfoLevel)
					}
				}
			case err := <-watch.Errors:
				{
					FeyLog.Errorf("Select watch error happened, errinfo: %v", err)
					return
				}
			}
		}
	}()

	select {}

}
