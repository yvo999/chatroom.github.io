package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"redis5/chatroom/common/message"
	"redis5/chatroom/common/utils"
	"redis5/chatroom/server/model"
)

type UserProcess struct {
	//字段
	Conn net.Conn
	//增加一个字段表明该Conn是哪个用户的
	UserId int
}

//编写通知所有在线的用户的方法
//userId通知其他在线用户,我上线
func (this *UserProcess) NotifyOtherOnlineUser(userId int) {
	//遍历UserMgr，然后一个一个的发送
	for id, up := range userMgr.onlineUsers {
		//过滤掉自己
		if id == userId {
			continue
		}
		//开始通知(单独方法)
		up.NotifyMeOnline(userId)
	}
}
func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("notifyUSM j.M", err)
		return
	}
	//将序列化后的NTFUSM赋给mes.data
	mes.Data = string(data)
	//对mes再次序列化.准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Printf("notifyML err: %v\n", err)
		return
	}
	//发送,创建tranfer实例发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("ntfyMO tf.W err: %v\n", err)
		return
	}
}

//编写一个函数serverProcessLogin函数，专门处理登陆请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	fmt.Println("mes:", mes)
	//核心代码
	//1.先从mes中取出mes.Data,并反序列化成loginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Printf("login json.unmarshal err: %v\n", err)
		return
	}
	//1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//2.再声明一个loginResMes,并完成赋值
	var loginResMes message.LoginResMes
	//到redis数据库完成验证
	_, err = model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
	} else {
		loginResMes.Code = 200
		//这里因为用户登陆成功，我们就该把登陆成功放入UserMgr中
		//将登陆成功的用户Id赋给THIS
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		//通知其他的在线用户我上线了
		this.NotifyOtherOnlineUser(loginMes.UserId)
		//将当前登陆的用户id保存在loginResMes.UsersId
		//遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println("登陆成功...")
	}
	//如果用户的id为100，密码为123456则合法，否则不合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	//合法
	//loginResMes.Code = 200
	//} else {
	//不合法
	//	loginResMes.Code = 500
	//	loginResMes.Error = "该用户不存在请注册再使用。。。"
	//}
	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Printf("loginResMes json.marshal err: %v\n", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.将resMes进行序列化准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Printf("resMes j.Merr: %v\n", err)
		return
	}
	//6.发送将其封装到Write
	//因为使用了分层模式（MVC），我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("writePkg err: %v\n", err)
	}
	return
}

//注册
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//1.先从mes中取出mes.Data,并反序列化成registerMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Printf("login json.unmarshal err: %v\n", err)
		return
	}
	//1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	//2.再声明一个registerResMes,并完成赋值
	var registerResMes message.RegisterResMes
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
	}
	//3.将registerResMes序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Printf("loginResMes json.marshal err: %v\n", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.将resMes进行序列化准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Printf("resMes j.Merr: %v\n", err)
		return
	}
	//6.发送将其封装到Write
	//因为使用了分层模式（MVC），我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Printf("writePkg err: %v\n", err)
	}
	return
}
