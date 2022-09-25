package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"redis5/chatroom/common/message"
	"redis5/chatroom/common/utils"
)

//显示登陆成功后的界面
func ShowMenu() {
	fmt.Println("--------恭喜登陆成功--------")
	fmt.Println("--------1.显示在线用户列表--------")
	fmt.Println("--------2.发送消息--------")
	fmt.Println("--------3.信息列表--------")
	fmt.Println("--------4.退出系统--------")
	fmt.Println("--------请选择(1~4):--------")
	var key int
	var content string
	//因为我们总会用到SmsProcess实例，所以定义在swtich外部
	smsPrcoess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		//fmt.Println("发送消息")
		fmt.Println("请输入您想群发的消息:")
		fmt.Scanf("%s\n", &content)
		err := smsPrcoess.SendGroupMes(content)
		if err != nil {
			fmt.Printf("sms.Send err: %v\n", err)
		}
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你选择退出系统...")
		os.Exit(0)
	default:
		fmt.Println("输入选项错误")
	}
}

//和服务器端保持通讯
func ProcessServerMes(conn net.Conn) {
	//创建一个transfer实例，不停读取服务器
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Printf("tf.RP err: %v\n", err)
			return
		}
		//如果读取到消息,又是下一步处理
		//fmt.Printf("mes: %v\n", mes)
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线
			//处理
			//1.取出NotifyUSM
			var notifyUserStatusMes message.NotifyUserStatusMes
			err = json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			if err != nil {
				fmt.Printf("nUSM j.M: %v\n", err)
			}
			//2.把这个用户的信息，状态保存到客户端的map
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType: //有人群发消息
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器端返回一个未知类型，暂时无法处理")
		}
	}
}
