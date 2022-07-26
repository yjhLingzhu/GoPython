package main

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

// init函数首先执行,初始化redis pool
func init() {
	pool = &redis.Pool{
		MaxIdle:     8,   // 最大空闲数
		MaxActive:   0,   // 最大连接数，0 代表没有限制
		IdleTimeout: 100, // 最大空闲时间 秒
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379", redis.DialDatabase(1))
		}, // 初始化连接代码，连接哪个ip
	}
}

func main() {
	// 从pool 取出一个连接
	conn := pool.Get()
	defer conn.Close()

	// 放进
	_, err := conn.Do("Set", "name", "ayun") // 把 "阿云" 存放进去redis中
	if err != nil {
		fmt.Println("放进去失败！")
		return
	}

	// 取出
	v, err := redis.String(conn.Do("Get", "name"))
	if err != nil {
		fmt.Println("取出失败！")
		return
	}
	fmt.Println("值为：", v)
}
