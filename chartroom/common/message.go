package common

const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
	RegisterMesType = "RegisterMes"
	ChatMesType     = "ChatMesType"
)

type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息数据
}

// 定义两个消息
type LoginMes struct {
	UserId   int    `json:"userId"`   // 用户id
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

type LoginResMes struct {
	Code  int    `json:"code"`  // 返回状态码
	Error string `json:"error"` // 错误信息
}

type RegisterMes struct {
	UserId   int    `json:"userId"`   // 用户id
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

type ChatResponse struct {
	Uid string
	Msg string
}
