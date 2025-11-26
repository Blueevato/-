package client_router

import (
	"encoding/json"
	"fmt"
	"gochat2/common/message"
	"gochat2/common/utils"
)

type SmsRouter struct {
}

// 发送群聊消息
func (this *SmsRouter) BroadCastSms(content string) (err error) {
	//封装mes
	var mes message.Message
	mes.Type = message.SmsMesType
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurentUser.UserId
	smsMes.UserStatus = CurentUser.UserStatus

	//序列化
	smsdata, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("smsdata marshal err", err)
		return
	}
	mes.Data = string(smsdata)
	senddata, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("broadcastsms marshal err", err)
		return
	}

	//发送给服务器
	tf := &utils.Transfer{
		Conn: CurentUser.Conn,
	}
	err = tf.WritePkg(senddata)
	if err != nil {
		fmt.Println("broadcastsms senddata err", err)
		return
	}
	return
}

// 显示即可
func OutoutBroadCastSms(mes *message.Message) {
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("OutoutBroadCastSms unmarshal err", err)
		return
	}

	//显示
	fmt.Printf("用户id:\t%d 对大家说:\t%s", smsMes.UserId, smsMes.Content)
	fmt.Println("\n")
}
