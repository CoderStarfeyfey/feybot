/**
 * @Description
 * @Author Xiongfei
 * @Date 2024/1/25
 **/
package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type HttpHook interface {
	BeforeRequestHook(req *http.Request)
	AfterRequestHook(response http.Response, err error)
}
type HttpHooks []HttpHook

// 通用框架中的http.client都会封装一层
type NewHttpClient struct {
	MaxRetryTimes int
	client        http.Client
}

func (c *NewHttpClient) do(req *http.Request) (*http.Response, error) {
	if c.MaxRetryTimes <= 0 {
		c.MaxRetryTimes = 3
	}
	var (
		resp *http.Response
		err  error
	)
	for i := 0; i < c.MaxRetryTimes; i++ {
		resp, err = c.client.Do(req)
		if err == nil {
			break
		}
	}
	return resp, err
}

// 封装http下载图片以及封装请求响应的结构体
func GetHttpRespJson(url string, RespStruct interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		FeyLog.Errorf("Error to get req,URL: %s,error:%v", url, err)
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(RespStruct)
	if err != nil {
		FeyLog.Errorf("Error to decoder json string,error:%v", err)
		return err
	}
	_, err = json.Marshal(RespStruct)
	if err != nil {
		FeyLog.Errorf("Error to parse struct to bytyes,error:%v", err)
		return err
	}
	return nil
}

func DownloadPicToDefaultDir(url string, path string) error {
	//还是先发送get请求在把内容存到path的文件中
	resp, err := http.Get(url)
	if err != nil {
		FeyLog.Errorf("Download Pic fail,url:%s,error:%v", url, err)
		return err
	}
	defer resp.Body.Close()
	//判断文件路径是否存在
	if _, err = os.Stat(path); err == nil {
		err = os.Remove(path)
		if err != nil {
			FeyLog.Errorf("Fail to delete fail,path:%s,Error:%v", path, err)
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		FeyLog.Debugf("filePath:%s create fail,Error:%v", path, err)
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		FeyLog.Debugf("fail to copy IO from url to file,Error:%v", err)
		return err
	}
	return nil
}
