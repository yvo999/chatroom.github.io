package main

import (
	"fmt"
	"io"
	"net"
	"redis5/chatroom/common/message"
	"redis5/chatroom/common/utils"
	processes "redis5/chatroom/server/process"
)

type Processor struct {
	Conn net.Conn
}

//编写一个serverProcessMes函数
//功能：根据消息类型不同而决定调用哪个函数处理
func (this *Processor) ServerProcessMes(mes *message.Message) (err error) {
	//是否能接收到从客户端发送的群发消息
	fmt.Println("mse=", mes)
	switch mes.Type {
	case message.LoginMesType:
		//处理登陆
		up := &processes.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
		if err != nil {
			fmt.Printf("up.SPLerr: %v\n", err)
		}
	case message.RegisterMesType:
		//处理注册
		up := &processes.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
		if err != nil {
			fmt.Printf("up.SPLerr: %v\n", err)
		}
	case message.SmsMesType:
		//创建smsprocess实例完成转发群聊消息
		smsProcess := &processes.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在 无法处理。。。")
	}
	return nil
}
func (this *Processor) Fasongprocess() (err error) {
	//循环的读客户端发送的信息
	for {
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出 服务器端退出...")
				return err
			} else {
				fmt.Printf("readPkg err: %v\n", err)
				return err
			}
		}
		//fmt.Printf("mes: %v\n", mes)
		err = this.ServerProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
