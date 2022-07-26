package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func ReadPkg(conn net.Conn) (mes Message, err error) {
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

	if mes.Type == ChatMesType {
		// 获取真正的data
		data := mes.Data
		// 将data反序列化成common.LoginMes
		var ChatData ChatResponse
		err = json.Unmarshal([]byte(data), &ChatData)
		if err != nil {
			fmt.Println("json.Unmarshal err", err)
			return
		}
		return
	}

	// 获取真正的data
	data := mes.Data
	// 将data反序列化成common.LoginMes
	var loginMes LoginMes
	err = json.Unmarshal([]byte(data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	return
}

func WritePkg(conn net.Conn, data []byte) (err error) {
	// 先发送一个长度给对方
	// 1 先把data的长度发送给服务器，防止丢包
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

	// 2 发送消息本身
	n, err = conn.Write(data)
	if err != nil || n != int(pkgLen) {
		fmt.Println("send data to server err ", err)
		return err
	}
	return
}
