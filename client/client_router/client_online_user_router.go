package client_router

import (
	"fmt"
	"gochat2/common/message"
	"gochat2/common/model"
)

// 客户端维护的onlinemap和当前用户
var (
	clientOnlineUsers map[int]*model.User = make(map[int]*model.User, 1024)
	CurentUser        model.CurentUser    //用户登录成功后初始化
)

// 客户端显示在线的用户
func showOnlineUsers() {
	fmt.Println("当前在线的用户：")
	for id, _ := range clientOnlineUsers {
		fmt.Println("用户id\t", id)
	}
}

// 处理函数
func updateUserStatus(notifyUserStatusMes *message.UserStatusNotifyMes) {

	user, ok := clientOnlineUsers[notifyUserStatusMes.UserId] //user已经是指针类型了
	if !ok {
		//不存在
		user = &model.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	clientOnlineUsers[notifyUserStatusMes.UserId] = user
	showOnlineUsers()
}
