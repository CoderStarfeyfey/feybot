package internal

import (
	"bufio"
	"fmt"
	"github.com/CoderStarfeyfey/feybot/internal/utils"
	"net"
	"strings"
)

// 自定义消息以及回复的结构体
type RequestStruct struct {
	Uname      string // 昵称
	Uuid       string // 用户唯一标识
	Groupname  string // 群昵称，取不到或没有的时候为""
	RequestTxt string // 请求字符串
}
type ReplyStruct struct {
	ReType int //回复类型 如果为1，就是普通消息此时的Retext为消息内容，如果为2，回复图片，此时为图片的路径,如果为3，回复文字+图片
	ReText string
}

type PluginFunc func(string, string) (*ReplyStruct, error)
type replyMsgStruct struct {
	ReType int
	ReText string
}

var PluginMap map[string]PluginFunc
var DefaultReplyMap = map[string]replyMsgStruct{
	"你的名字": {
		ReType: 1,
		ReText: "我的名字是feybot",
	},
	"你的年纪": {
		ReType: 1,
		ReText: "我的年龄是100岁",
	},
}

// 定义命令行参数注册结构体
type CommandHandler struct {
	Execute func(args []string) (string, error)
	Help    string
}

var CommandHandlers map[string]CommandHandler

func RegisterCommand(name string, handler CommandHandler) {
	CommandHandlers[name] = handler
}

// 这里需要写一个调试命令控制台服务端协程的入口函数
func DbgConsoleServiceHandle() {
	// 在本地6000端口监听TCP连接
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println("监听错误:", err)
		return
	}
	defer ln.Close()

	fmt.Println("服务端启动，正在监听...")

	for {
		// 接受新的连接
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("接受连接错误:", err)
			continue
		}

		// 处理连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// 读取客户端发送的消息
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("读取错误:", err)
			}
			break
		}
		message = strings.TrimSpace(message)
		fmt.Println("收到命令:", message)

		// 根据接收到的命令执行相应的操作
		response := processCommand(message)
		fmt.Sprintf(response)
		// 向客户端发送响应
		_, err = conn.Write([]byte(response + "<END>\n"))
		if err != nil {
			fmt.Println("发送响应失败:", err)
			break
		}
	}
}

// 解析调试命令并执行对应的回调函数
func processCommand(command string) string {
	//feybot tyqw -u uuid -h 打印该用户的历史缓存
	//feybot tyqw -u uuid -t 打印该已经使用的次数
	parts := strings.Fields(command)
	if len(parts) <= 2 {
		return "invalid args"
	}
	if parts[0] != "feybot" {
		return "invalid args"
	}
	featureName := parts[1]
	args := parts[2:]
	if hander, ok := CommandHandlers[featureName]; ok {
		if len(args) == 1 && args[0] == "?" {
			return hander.Help
		}
		res, err := hander.Execute(args)
		if err != nil {
			return "execute error"
		}
		return res
	}
	return fmt.Sprintf("%s do not register", featureName)

}

func init() {
	utils.FeyLog.Debug("create plugin map")
	PluginMap = make(map[string]PluginFunc)
	CommandHandlers = make(map[string]CommandHandler)
}
