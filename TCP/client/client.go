package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("dial err", err)
		return
	}
	// fmt.Println("conn=", conn)

	// 功能一：客户端可以发送单行数据，然后退出
	// 新建一个终端reader从终端读取东西
	reader := bufio.NewReader(os.Stdin) // os.Stdin 代表标准输入[即终端]

	// 从终端读取用户输入的数据发送给服务的
	for {
		fmt.Println("请输入信息：")
		line, err := reader.ReadString('\n') // 这里会阻塞，因为reader是标准输入
		if err != nil {
			fmt.Println("reader string err!", err)
		} else {
			// 将数据发送给服务器
			n, err := conn.Write([]byte(line))
			if err != nil {
				fmt.Println("Write err", err)
				continue
			}
			fmt.Println("n=", n)
		}
 
	}
}
