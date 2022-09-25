package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"redis5/chatroom/common/message"
	"redis5/chatroom/common/utils"
)

type SmsProcess struct {
}

//转发方法
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端map
	//将消息转发出去
	//取出mes的内容
	var smsTranferMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsTranferMes)
	if err != nil {
		fmt.Printf("smsMes j.U err: %v\n", err)
		return
	}
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Printf("SGM.mes.j.M err: %v\n", err)
		return
	}
	for id, up := range userMgr.onlineUsers {
		//这里，还需要过滤到自己，不要发给自己
		if id == smsTranferMes.UserId {
			continue
		}
		this.SendMesEachOnlineUser(data, up.Conn)
	}
}
func (this *SmsProcess) SendMesEachOnlineUser(data []byte, conn net.Conn) {
	//创建一个transfer实例发送
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Printf("SendMEOU.tf.WPerr: %v\n", err)
		return
	}
}
