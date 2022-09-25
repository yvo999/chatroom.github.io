package processes

import (
	"encoding/json"
	"fmt"
	"redis5/chatroom/common/message"
)

func outputGroupMes(mes *message.Message) { //这个地方mes一定是SmsMes
	//显示即可
	//1.反序列化
	var smsTranferMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsTranferMes)
	if err != nil {
		fmt.Printf("ouputGM.sTM.j.S err: %v\n", err)
		return
	}
	//显示信息
	info := fmt.Sprintf("用户id:\t%d对大家说:\t%s", smsTranferMes.UserId, smsTranferMes.Content)
	fmt.Println(info)
	fmt.Println()
}
