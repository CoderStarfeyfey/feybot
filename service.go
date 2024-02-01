package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
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

		// 向客户端发送响应
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("发送响应失败:", err)
			break
		}
	}
}

func processCommand(command string) string {
	// 在这里根据接收到的命令进行处理
	// 示例：简单地返回一个确认信息
	return "命令已接收: " + command
}
