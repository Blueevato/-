package main

import (
	"fmt"
	"gochat2/common/message"
	"gochat2/common/utils"
	"gochat2/server/router"
	_ "gochat2/server/router"
	"io"
	"net"
)

type ProcessRouter struct {
	Conn net.Conn
}

// 总的处理消息路由
func (this *ProcessRouter) Process_router(mes *message.Message) (err error) {
	//test
	fmt.Println("mes=", mes)

	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		ur := &router.UserRouter{
			Conn: this.Conn,
		}
		err = ur.Process_login(mes)
	case message.RegisterMesType:
		//处理注册
		ur := &router.UserRouter{
			Conn: this.Conn,
		}
		err = ur.Process_register(mes)
	case message.SmsMesType:
		//处理群发or单发消息
		//TODO
		sr := &router.SmsRouter{}
		err = sr.Process_smsBroadcast(mes)
	default:
		fmt.Println("路由err...")
	}

	return
}
func (this *ProcessRouter) MainProcess() (err error) {
	//read
	for {
		//读数据包
		//创建一个Transfer实例
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		var mess message.Message
		mess, err = tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端断开连接...")
				break
			}
			fmt.Println("read pkg err", err)
			continue
		}
		fmt.Println("mess=", mess)

		//调用路由
		err = this.Process_router(&mess)
		if err != nil {
			fmt.Println("process_router err", err)
			return
		}

	}
	return
}
