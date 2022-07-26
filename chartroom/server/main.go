package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_python/chartroom/common"
	"io"
	"net"
)

// 这个方法被抽到了common下的utils里
func readPkg(conn net.Conn) (mes common.Message, err error) {
	buf := make([]byte, 8096)
	// 在conn没有关闭的情况下才会阻塞，如果客户端关闭了conn则不会阻塞
	_, err = conn.Read(buf[:4]) // 因为首先传过来的是长度
	if err != nil {
		// fmt.Println("conn.Read buf[:4] err ", err)
		return
	}
	// 根据buf[:4]转成一个unit32类型
	pkgLen := binary.BigEndian.Uint32(buf[:4])
	// 根据pkgLen读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if err != nil || n != int(pkgLen) {
		fmt.Println("conn.Read buf[:pkgLen] err ", err)
		return
	}

	// 把pkgLen反序列化成 -> common.Message
	err = json.Unmarshal(buf[:pkgLen], &mes) // 一定要加 &
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	// 获取真正的data
	data := mes.Data
	// 将data反序列化成common.LoginMes
	var loginMes common.LoginMes
	err = json.Unmarshal([]byte(data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	return
}

// 这个方法被抽到了common下的utils里
func writePkg(conn net.Conn, data []byte) (err error) {
	// 先发送一个长度给对方
	// 1 先把data的长度发送给服务器，防止丢包
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	fmt.Println("buf:", buf)
	// 发送长度
	n, err := conn.Write(buf[:])
	if err != nil || n != 4 {
		fmt.Println("send data len to server err ", err)
		return err
	}

	// 2 发送消息本身
	n, err = conn.Write(data)
	if err != nil || n != int(pkgLen) {
		fmt.Println("send data to server err ", err)
		return err
	}
	return
}

// 编写一个函数serverProcessLogin函数， 专门处理登录请求
func serverProcessLogin(conn net.Conn, mes *common.Message) (err error) {
	// 先从mes 中取出 mes.Data， 并直接反序列化成LoginMes
	var loginMes common.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("获取data失败")
		return
	}

	// 先声明一个resMes
	var resMes common.Message
	resMes.Type = common.LoginResMesType

	// 再声明一个返回数据类型
	var loginResMes common.LoginResMes
	// 如果id=100, pwd=yjh，合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "yjh" {
		// 登录成功
		fmt.Println("登录成功")
		loginResMes.Code = 200
	} else {
		// 不合法
		fmt.Println("登录失败")
		loginResMes.Code = 400
		loginResMes.Error = "用户不存在请先注册！"
	}
	// 序列化返回信息
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("解析登录返回数据时失败 loginResMes ....")
		return
	}
	resMes.Data = string(data)

	// 对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("解析登录返回数据时失败 resMes ....")
		return
	}
	// 发送response回去
	err = common.WritePkg(conn, data)
	if err != nil {
		fmt.Println("发送response err ", err)
	}
	return
}

// 编写一个函数ChatProcess，专门处理聊天的
func ChatProcess(conn net.Conn, mes *common.Message) (err error) {
	// 返回的Data
	var chatData, resChatData common.ChatResponse
	err = json.Unmarshal([]byte(mes.Data), &chatData)
	if err != nil {
		fmt.Println("获取聊天数据失败")
		return
	}
	// 先声明一个resMes
	var resMes common.Message
	resMes.Type = common.ChatMesType
	fmt.Print(chatData.Msg)

	if chatData.Msg == "我来了\r\n" {
		resChatData.Uid = "xiao Y"
		resChatData.Msg = "欢迎来到yjh小天地聊天室~~~"
	} else {
		resChatData.Uid = "xiao Y"
		resChatData.Msg = "哎呀，小Y听不懂你在说什么~"
	}

	// 对resChatData进行序列化
	data, err := json.Marshal(resChatData)
	if err != nil {
		fmt.Println("解析聊天返回数据时失败 resMes ....")
		return
	}
	resMes.Data = string(data)
	// 对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("解析聊天返回数据时失败 resMes ....")
		return
	}
	// 发送response回去
	err = common.WritePkg(conn, data)
	if err != nil {
		fmt.Println("发送response err ", err)
		return
	}
	return
}

// 编写一个ServerProcessMes 函数
// 功能：根据客户的消息种类，决定调用哪个函数来处理
func serverProcessMes(conn net.Conn, mes *common.Message) (err error) {
	switch mes.Type {
	case common.LoginMesType:
		// 处理登录
		err = serverProcessLogin(conn, mes)
		if err != nil {
			fmt.Println("登录失败")
		}
	case common.RegisterMesType:
		// 处理注册
	case common.ChatMesType:
		// 开始聊天
		err = ChatProcess(conn, mes)
		if err != nil {
			fmt.Println("服务器异常！聊天室崩溃了啦....")
		}
	default:
		fmt.Println("消息类型不存在，无法处理....")
	}
	return
}

// 处理客户端的通讯
func process(conn net.Conn) {
	defer conn.Close()
	// 获取客户端发送的信息
	for {
		fmt.Println("读取客户端数据....")
		// 这里将读取数据包直接封装成一个readPkg()函数
		// mes, err := readPkg(conn)
		mes, err := common.ReadPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出了，服务端也退出")
				return
			} else {
				fmt.Println("readPkg err ", err)
				fmt.Println("客户端退出了，服务端也退出")
				return
			}
		}
		err = serverProcessMes(conn, &mes)
		if err != nil {
			fmt.Println("处理客户消息出错了....")
			return
		}

	}
}

func main() {
	fmt.Println("服务器开始监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()

	if err != nil {
		fmt.Println("listen err ", err)
		return
	}

	for {
		fmt.Println("等待客户端来连接....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("获取连接失败！", err)
		}

		// 连接成功，启动一个协程跟客户端保持通讯
		go process(conn)
	}
}
