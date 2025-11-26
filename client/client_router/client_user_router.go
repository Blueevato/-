package client_router

import (
	"encoding/json"
	"fmt"
	"gochat2/common/message"
	"gochat2/common/model"
	"gochat2/common/utils"
	"net"
	"os"
)

type ClientUserRouter struct {
	//	Conn net.Conn
}

func (this *ClientUserRouter) Login(userId int, userPwd string) (err error) {

	//fmt.Printf("userId =%d userPwd=%s\n", userId, userPwd)
	//return nil

	//1 链接到服务器
	conn, err := net.Dial("tcp", "localhost:"+"8889")
	if err != nil {
		fmt.Println("dial err")
		return
	}
	defer conn.Close()

	//2 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	//3 创建LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	data, err := json.Marshal(loginMes) //序列化
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	mes.Data = string(data)

	mesdata, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(mesdata)
	if err != nil {
		fmt.Println("登录 write pkg err", err)
		return
	}

	//接收服务端的消息
	var resmes message.Message
	resmes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("read pkg err", err)
		return
	}

	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(resmes.Data), &loginResMes)
	if err != nil {
		fmt.Println("resmes.Data unmarshal err", err)
		return
	}

	if loginResMes.CodeId == 200 {
		fmt.Println("login succ!")

		//初始化CurentUser
		CurentUser.Conn = conn
		CurentUser.UserId = userId
		CurentUser.UserStatus = message.UserOnline

		cr := &ClientRouter{
			Conn:     conn,
			UserName: loginMes.UserName,
		}

		fmt.Println("当前在线用户如下：")
		for _, v := range loginResMes.UserIds {
			if v == userId {
				continue
			}
			fmt.Println("用户id：\t", v)
			//完成客户端的clientOnlineUsers的初始化
			user := &model.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			clientOnlineUsers[v] = user
		}
		fmt.Println("\n\n")

		//go一个协程和服务器保持通讯
		//如果服务器有数据推送给客户端，则接收并显示
		go cr.KeepServerMes()

		//TODO 进入二级菜单
		//1 显示登录成功后的二级菜单
		cr.ShowMenu()

	} else { //登录失败
		fmt.Println("login err!", loginResMes.ErrorMsg)
	}

	return

}

func (this *ClientUserRouter) Register(userId int, userPwd, userName string) (err error) {
	//1 链接到服务器
	conn, err := net.Dial("tcp", "localhost:"+"8889")
	if err != nil {
		fmt.Println("dial err")
		return
	}
	defer conn.Close()

	//2 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	//3 创建RegisterMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	data, err := json.Marshal(registerMes) //序列化
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	mes.Data = string(data)

	mesdata, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(mesdata)
	if err != nil {
		fmt.Println("注册 write pkg err", err)
		return
	}

	//接收服务端的消息
	var resmes message.Message
	resmes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("read pkg err", err)
		return
	}
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(resmes.Data), &registerResMes)
	if registerResMes.CodeId == 200 {
		fmt.Println("register succ!")

		//TODO if

	} else { //注册失败
		fmt.Println("register err!", registerResMes.ErrorMsg)
	}

	conn.Close()
	os.Exit(0)

	return

}
