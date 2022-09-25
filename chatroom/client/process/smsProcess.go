package processes

import (
	"encoding/json"
	"fmt"
	"redis5/chatroom/common/message"
	"redis5/chatroom/common/utils"
)

type SmsProcess struct {
}

//发送群聊的消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个mes
	var mes message.Message
	mes.Type = message.SmsMesType
	//2.创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content //内容
	smsMes.User.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus
	//3.序列化smsmes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Printf("smsmes j.S err: %v\n", err)
		return
	}
	mes.Data = string(data)
	//4.对mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Printf("mes j.M err: %v\n", err)
		return
	}
	//5.将序列化后的mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("sms tf.WP err: %v\n", err)
		return
	}
	return
}
