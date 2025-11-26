package client_router

import (
	"encoding/json"
	"fmt"
	"gochat2/common/message"
	"gochat2/common/utils"
	"net"
	"os"
)

type ClientRouter struct {
	Conn     net.Conn
	UserName string
}

// 登陆成功后的界面
func (this *ClientRouter) ShowMenu() {
	key := 0
	var content string
	smsRouter := &SmsRouter{}
	for {
		//TODO name
		fmt.Println("--user " + this.UserName + " ,succ login...--")
		fmt.Println("--\t\t\t 1 显示在线用户--")
		fmt.Println("--\t\t\t 2 发送消息--")
		fmt.Println("--\t\t\t 3 消息列表--")
		fmt.Println("--\t\t\t 4 退出--")
		fmt.Println("--\t\t\t 请选择(1-4):--")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("--显示在线用户--")
			//TODO 管道实现广播
			showOnlineUsers()
		case 2:
			fmt.Println("--发送消息--")
			fmt.Scanf("%s\n", &content)
			smsRouter.BroadCastSms(content)
		case 3:
			fmt.Println("--消息列表--")
		case 4:
			fmt.Println("--exit--")
			// 在退出前关闭连接
			if this.Conn != nil {
				this.Conn.Close()
			}
			os.Exit(0)
		default:
			fmt.Println("--input err--")
		}
	}
}

// 和服务器保持通讯
func (this *ClientRouter) KeepServerMes() {

	//创建一个transfer实例，不停读取服务器发来的消息
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发来的消息...")
		mes, err := tf.ReadPkg()
		if err != nil {

			fmt.Println("readpkg err", err)
			return
		}

		//读取到消息，进行路由
		switch mes.Type {
		case message.UserStatusNotifyMesType:
			//有人上线
			var notifuUserStatusMes message.UserStatusNotifyMes
			err = json.Unmarshal([]byte(mes.Data), &notifuUserStatusMes)
			if err != nil {
				fmt.Println("UserStatusNotifyMes unmarshal err", err)
				continue
			}
			updateUserStatus(&notifuUserStatusMes)
		case message.SmsMesType:
			//有人群发消息
			OutoutBroadCastSms(&mes)
		default:
			fmt.Println("default Type")
		}

		//fmt.Println("mes=", mes)
	}

}
