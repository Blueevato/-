package router

import (
	"encoding/json"
	"fmt"
	"gochat2/common/message"
	"gochat2/common/utils"
	"net"
)

// 处理和短消息有关的请求
type SmsRouter struct {
	//Conn net.Conn
}

func (this *SmsRouter) Process_smsBroadcast(mes *message.Message) (err error) {
	//取出信息
	var smsMes message.SmsMes
	err = json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("smsMes unmarshal err", err)
		return
	}

	//序列化数据
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("smsMes marshal err", err)
		return
	}

	//遍历服务器的onlinemap，将消息群发
	for id, ur := range onlineUserRouter.onlineUsers {
		//过滤自己
		if id == smsMes.UserId {
			continue
		}

		this.Process_smsBroadcastToOne(data, ur.Conn)
	}
	return
}

func (this *SmsRouter) Process_smsBroadcastToOne(content []byte, conn net.Conn) {
	//发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(content)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
		return
	}
}
