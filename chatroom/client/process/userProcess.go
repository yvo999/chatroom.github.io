package processes

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"redis5/chatroom/common/message"
	"redis5/chatroom/common/utils"
)

type UserPorces struct {
}

//写一个函数完成登陆效验
func (this *UserPorces) Login(userId int, userPwd string) (err error) {
	//下一步开始定协议
	//fmt.Printf("userID=%d userPWD=%s\n", userId, userPWD)
	//return nil
	//1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Printf("net.dial err: %v\n", err)
		return
	}
	defer conn.Close()
	//2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	//3.创建loginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	//4.将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Printf("loginMes json.marshal err: %v\n", err)
		return
	}
	//把data赋值给mes
	mes.Data = string(data)
	//6.将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Printf("mes.data jm err: %v\n", err)
		return
	}
	//7.到这个时候data就是要发送的消息
	//7.1先把data长度发送给服务器
	//先获取到data的长度->转成一个长度的byte切片
	pkgLen := uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Printf("conn.write err: %v\n", err)
		return
	}
	//fmt.Println("客户端发送消息的长度一致且成功")
	fmt.Printf("消息的长度为%d,内容为%s\n", len(data[0:4]), string(data))
	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("conn.write(data)err: %v\n", err)
		return
	}
	fmt.Printf("mes.Data: %v\n", mes.Data)
	//休眠10秒
	//time.Sleep(10 * time.Second)
	//fmt.Println("休眠了10秒...")
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Printf("mes u.Rerr: %v\n", err)
		return
	}
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Printf("loginResMes j.U err: %v\n", err)
		return
	} else if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline
		//显示当前在线用户的列表
		fmt.Println("当前在线用户列表如下")
		for _, v := range loginResMes.UsersId {
			//不显示自己的id
			if v == userId {
				continue
			}
			fmt.Printf("当前在线用户id: %v\n", v)
			//完成 客户端的onlineUsers完成初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		//这里我们还需要启动一个协程
		//该协程保持与服务器的通讯，如果服务器有数据推送给客户端
		//则接收并显示在客户端
		go ProcessServerMes(conn)
		//1.显示登陆成功的菜单[循环]
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}

//请求注册方法
func (this *UserPorces) Register(userId int, userPwd string, userName string) (err error) {
	//链接
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Printf("Register net.Dial err: %v\n", err)
		return
	}
	defer conn.Close()
	//2.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	//3.创建loginMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	//4.将loginMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Printf("register j.M err: %v\n", err)
		return
	}
	//把data赋值给mes
	mes.Data = string(data)
	//6.将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Printf("register mes j.M err: %v\n", err)
		return
	}
	tf := &utils.Transfer{
		Conn: conn,
	}
	//发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("register tf.WP err: %v\n", err)
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Printf("register mes tf.RP err: %v\n", err)
		return
	}
	//将mes的data部分反序列化成RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登陆")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}
