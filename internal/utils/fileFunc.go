/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/25
 **/
package utils

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func RandomLessThan(max int) int {
	if max <= 1 {
		// 如果 max 小于或等于 1，则无法生成符合条件的随机数
		return 0
	}

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成 [0, max-1) 的随机数然后加 1
	return rand.Intn(max - 1)
}

func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
func StringSubContains(slice []string, element string) bool {
	for _, v := range slice {
		if strings.Contains(element, v) {
			return true
		}
	}
	return false
}
func GenPicTempPath() (string, error) {
	//先检临时文件名是否存在
	if _, err := os.Stat(TmpDir); os.IsNotExist(err) {
		err = os.MkdirAll(TmpDir, 0755)
		if err != nil {
			FeyLog.Errorf("fail to create TmpDir:%s,Error:%v", TmpDir, err)
			return "", err
		}
	}
	currentTime := time.Now().Format("20060102")

	return fmt.Sprintf("%s/%s", TmpDir, currentTime), nil

}
func MapContainValue[K comparable, V comparable](m map[K]V, element V) (bool, K) {
	for key, value := range m {
		if value == element {
			return true, key
		}
	}
	var zeroK K
	return false, zeroK
}

func containsValue[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}
	return false
}
