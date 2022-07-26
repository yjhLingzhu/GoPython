package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

// 要求是客户端A和客户端B通信

func handleConnection(conn net.Conn, clientInfos chan string) {
	defer conn.Close()
	for {
		// 创建一个切片
		buf := make([]byte, 1024)
		fmt.Printf("服务器在等待%s客户端发信息！\n", conn.RemoteAddr().String())

		// 创建一个变量记录当前客户端的ip，防止连接断开后不知道从redis中该移除哪个IP
		currentIP := conn.RemoteAddr().String()

		// 1. 等待客户端发送的信息
		// 2. 如果客户端没有write[发送], 那么协程会阻塞在这里直到超时
		n, err := conn.Read(buf) // 从conn读取	这里会阻塞
		if err != nil {
			// 客户端断开连接时，从redis中移除相应的ip
			ctx := context.Background()
			rdb.SRem(ctx, "user1", currentIP)
			fmt.Println("服务器端read err", err)
			break
		}
		// 3. 显示信息到服务器终端
		fmt.Print("信息是：", string(buf[:n]))

		// 将信息写到公共管道中
		clientInfos <- string(buf[:n])
	}
}

// 启一个协程，一直监听公共管道，一旦发现有信息进入立马获取，并遍历redis中的全部ip，然后发送给客户端
func sendMsgToClient(clientInfos chan string) {
	for {
		v, ok := <-clientInfos
		if !ok {
			fmt.Println("没有消息了...")
			break
		}
		// 获取redis中所有ip
		ctx := context.Background()
		es, _ := rdb.SMembers(ctx, "user1").Result()
		// fmt.Println("结果：", es)
		for _, ip := range es {
			// 这里向各个客户端发送消息
			fmt.Println("客户：", ip, "消息：", v)
		}
	}
}

func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",  // no password set
		DB:       1,   // use default DB
		PoolSize: 100, // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	return err
}

func main() {
	fmt.Println("服务器开始监听....")
	// 监听tcp协议的8888端口
	listen, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("listen err!", err)
		return
	}
	defer listen.Close()

	// 新建一个管道，用来接受各个客户端发送过来的信息
	var clientInfos = make(chan string, 1000)
	go sendMsgToClient(clientInfos) // 向其他客户端发送消息

	// 初始化redis
	ctx := context.Background()
	if err := initClient(); err != nil {
		fmt.Println("redis数据库连接失败！")
		return
	}

	// 循环等待客户端来连接
	for {
		// 等待客户端连接，阻塞
		fmt.Println("等待客户端来连接...")
		conn, err := listen.Accept() // 阻塞
		if err != nil {
			fmt.Println("accept err!", err)
			continue
		} else {
			// 将ip写到redis 去重
			err := rdb.SAdd(ctx, "user1", conn.RemoteAddr().String()).Err()
			if err != nil {
				panic(err)
			}

			fmt.Printf("conn=%v, addr=%s\n", conn, conn.RemoteAddr().String())
		}
		// 这里开协程为客户端服务
		go handleConnection(conn, clientInfos)

	}

}
