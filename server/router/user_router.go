package router

import (
	"encoding/json"
	"fmt"
	"gochat2/common/message"
	"gochat2/common/model"
	"gochat2/common/utils"
	"gochat2/server/dao"
	"net"
)

type UserRouter struct {
	Conn   net.Conn
	UserId int //该Conn是哪个用户的
}

// 通知所有在线用户的方法
func (this *UserRouter) NotifyOnlineOthers(userId int) {
	//遍历
	for id, ur := range onlineUserRouter.onlineUsers {
		//过滤自己
		if id == userId {
			continue
		}

		//开始通知
		ur.NotifyOnlineOne(userId)

	}
}

// 具体的通知方法
func (this *UserRouter) NotifyOnlineOne(userId int) {
	//封装消息
	var mes message.Message
	mes.Type = message.UserStatusNotifyMesType
	var notifyMes message.UserStatusNotifyMes
	notifyMes.UserId = userId
	notifyMes.Status = message.UserOnline

	//序列化 ,一层一层封装，先封装notifyMes
	notifydata, err := json.Marshal(notifyMes)
	if err != nil {
		fmt.Println("notityMes.Data marshal err", err)
		return
	}
	mes.Data = string(notifydata) //序列化的数据赋值给Data字段
	senddata, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("notityMes marshal err", err)
		return
	}

	//发送通知消息
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(senddata)
	if err != nil {
		fmt.Println("notityMes WritePkg err", err)
		return
	}
}

// 处理和用户有关的请求
// 处理登录的核心代码
func (this *UserRouter) Process_login(mes *message.Message) (err error) {

	//1从mes取出 mes.Data,并反序列化还原
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("unmarshal err", err)
		return
	}

	//2 判断用户信息
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var loginResMes message.LoginResMes

	//去redis验证
	user, err := dao.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		fmt.Println("dao login fail", err)
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.CodeId = 500
			loginResMes.ErrorMsg = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.CodeId = 501
			loginResMes.ErrorMsg = err.Error()
		} else {
			loginResMes.CodeId = 502
			loginResMes.ErrorMsg = "服务器内部错误..."
		}
	} else {
		loginResMes.CodeId = 200
		loginResMes.ErrorMsg = ""
		fmt.Println("登录成功...", user)

		//登录成功，将用户添加到onlinemap,并通知其他用户
		this.UserId = loginMes.UserId
		onlineUserRouter.AddOnlineUser(this)
		this.NotifyOnlineOthers(loginMes.UserId)

		//遍历onlinemap,写入loginResMes的UserIds数组
		for id, _ := range onlineUserRouter.onlineUsers {
			loginResMes.UserIds = append(loginResMes.UserIds, id)
		}

	}

	//test id 100 pwd 123
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123" {
	//	//合法
	//	loginResMes.CodeId = 200
	//	loginResMes.ErrorMsg = ""
	//} else {
	//	//不合法
	//	loginResMes.CodeId = 500
	//	loginResMes.ErrorMsg = "登录失败,用户不存在..."
	//}

	//3 将 loginResMes 序列化,并赋值给resMes
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	resMes.Data = string(data)

	//4 序列化resMes并发送
	resdata, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}

	//5 发送给客户端,封装到writePkg()
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(resdata)

	return
}

func (this *UserRouter) Process_register(mes *message.Message) (err error) {

	//1从mes取出 mes.Data,并反序列化还原
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("unmarshal err", err)
		return
	}

	//2 提取用户信息
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//去redis写入
	err = dao.MyUserDao.Register(registerMes.User)
	if err != nil {
		fmt.Println("dao register fail", err)
		if err == model.ERROR_USER_EXISTS {
			registerResMes.CodeId = 504
			registerResMes.ErrorMsg = err.Error()
		} else {
			registerResMes.CodeId = 505
			registerResMes.ErrorMsg = "注册出现未知错误..."
		}
	} else {
		registerResMes.CodeId = 200
		registerResMes.ErrorMsg = ""
		fmt.Println("注册成功...")
	}

	//3 将 loginResMes 序列化,并赋值给resMes
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	resMes.Data = string(data)

	//4 序列化resMes并发送
	resdata, err := json.Marshal(resMes)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}

	//5 发送给客户端,封装到writePkg()
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(resdata)

	return
}
