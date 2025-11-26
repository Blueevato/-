package message

import "gochat2/common/model"

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	UserStatusNotifyMesType = "UserStatusNotifyMes"
	SmsMesType              = "SmsMes"
)

const (
	UserOnline     = iota //在线
	UserOffline           //离线
	UserBusyStatus        //繁忙

)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息的内容
}

type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	CodeId   int    `json:"codeId"` //返回的状态码
	UserIds  []int  //保存用户id的切片，用于功能：显示在线用户
	ErrorMsg string `json:"errorMsg"` //返回的错误信息
}

type RegisterMes struct {
	User model.User `json:"user"`
}

type RegisterResMes struct {
	CodeId   int    `json:"codeId"`   //返回的状态码
	ErrorMsg string `json:"errorMsg"` //返回的错误信息
}

// 服务器推送用户状态变化的消息
type UserStatusNotifyMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

// 短消息
type SmsMes struct {
	Content    string `json:"content"` //消息内容
	model.User        //匿名结构体
}
