package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_python/chartroom/common"
	"net"
	"os"
)

// 登录聊天室
func login(userId int, userPwd string) error {
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err ", err)
		return err
	}
	defer conn.Close()

	// 2. 发送消息给服务器
	var mes common.Message
	mes.Type = common.LoginMesType
	// 3. 创建LoginMes 结构体
	var loginMes common.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4. 将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json marshal loginMes err ", err)
		return err
	}
	// 5. 把data赋给mes.Data字段
	mes.Data = string(data)

	// 6. 将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal mes err ", err)
	}
	// 7. 将数据发送给服务器
	// 7.1 先把data的长度发送给服务器，防止丢包
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	// fmt.Println("buf:", buf)
	// 发送长度
	n, err := conn.Write(buf[:])
	if err != nil || n != 4 {
		fmt.Println("send data len to server err ", err)
		return err
	}

	// 7.2 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("send data to server err ", err)
		return err
	}

	// 8. 处理服务器端返回的消息
	// fmt.Println("客户端发送的长度=", len(data))
	// time.Sleep(time.Second * 10)
	mes, err = common.ReadPkg(conn)
	if err != nil {
		fmt.Println("读取服务器返回的数据失败！")
		return err
	}
	// 8.1 将mes的Data反序列化成LoginResMes
	var loginResMes common.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("解析mes.Data数据失败！")
		return err
	}
	// 8.2 判断状态码
	if loginResMes.Code != 200 {
		fmt.Println(loginResMes.Error)
		return nil
	} else {
		fmt.Println("登录成功！")
		// 新建一个终端reader从终端读取数据
		reader := bufio.NewReader(os.Stdin)
		// 这里for循环可以一直实现和服务器进行交互
		for {
			fmt.Println("请输入消息：")
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("输入数据有误！")
				continue
			}
			// 8.3 发送数据
			var mesData common.ChatResponse
			mesData.Uid = "yjh"
			mesData.Msg = line
			// 对mesData进行序列化
			msgData, err := json.Marshal(mesData)
			if err != nil {
				fmt.Println("发送数据前的序列化数据失败！")
				continue
			}
			// 先声明一个Message
			var Mes common.Message
			Mes.Type = common.ChatMesType
			Mes.Data = string(msgData)
			// 对Mes进行序列化，准备发送
			dataChat, err := json.Marshal(Mes)
			if err != nil {
				fmt.Println("发送数据前的序列化数据失败！")
				continue
			}
			err = common.WritePkg(conn, dataChat)
			if err != nil {
				fmt.Println("发送数据失败！")
				continue
			}
			// 8.4 读取数据
			mes, err := common.ReadPkg(conn)
			if err != nil {
				fmt.Println("读取服务器返回的数据失败！")
				return err
			}

			// 获取真正的data
			data := mes.Data
			// 将data反序列化成common.ChatResponse
			var ChatData common.ChatResponse
			err = json.Unmarshal([]byte(data), &ChatData)
			if err != nil {
				fmt.Println("json.Unmarshal err", err)
				continue
			}
			fmt.Println("服务器：", ChatData.Msg)
		}
	}
	// return nil
}
