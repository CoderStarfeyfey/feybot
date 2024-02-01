package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6000")
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("请输入命令: ")
		if scanner.Scan() {
			text := scanner.Text()
			_, err := conn.Write([]byte(text + "\n"))
			if err != nil {
				fmt.Println("发送命令失败:", err)
				return
			}

			// 设置10秒的读取超时
			err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				fmt.Println("设置读取超时失败:", err)
				return
			}

			// 等待响应
			response, err := readResponse(conn)
			if err != nil {
				fmt.Println("读取响应失败:", err)
				return
			}

			fmt.Println("收到响应:", response)
		}
	}
}

// readResponse 读取直到超时或遇到特定条件
func readResponse(conn net.Conn) (string, error) {
	var sb strings.Builder
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return sb.String(), err
		}
		sb.WriteString(line)
		if strings.HasSuffix(line, "<END>\n") {
			break
		}
	}
	return sb.String(), nil
}
