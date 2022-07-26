package main

import (
	"fmt"
)

func main() {
	// 定义一个接受用户输入的变量
	var key int
	// 是否显示菜单
	loop := true

	// 用户密码
	var userId int
	var userPwd string

	for loop {
		fmt.Println("--------------------欢迎登录多人聊天系统--------------------")
		fmt.Println("\t\t\t 1. 登录聊天室\t\t\t")
		fmt.Println("\t\t\t 2. 用户注册\t\t\t")
		fmt.Println("\t\t\t 3. 退出系统\t\t\t")
		fmt.Println("\t\t\t 请选择(1-3)\t\t\t")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			err := login(userId, userPwd)
			if err != nil {
				fmt.Println("登录失败！")
			} else {
				loop = false
			}

		case 2:
			fmt.Println("登录聊天室")
			loop = false
		case 3:
			fmt.Println("登录聊天室")
			loop = false
		default:
			fmt.Println("输入有误，请重新输入！")
		}
	}
}
