package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

//定义几个用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息的类型
	Data string `json:"data"` //消息
}

//定义两个消息后面需要再加
type LoginMes struct {
	UserId   int    `json:"userId"`   //用户ID
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}
type LoginResMes struct {
	Code    int    `json:"code"` //返回状态码 500为用户未注册 200为登陆成功
	UsersId []int  //增加字段，保存用户id的切片
	Error   string `json:"error"` //返回错误信息
}
type RegisterMes struct {
	User User `json:"user"`
}
type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码400为用户名已占有，200表示成功
	Error string `json:"error"` //返回错误信息
}

//为了配合服务器端推送用户状态变化的消息类型
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` //用户id
	Status int `json:"status"` //用户状态
}

//增加一个SmsMes 发送的消息
type SmsMes struct {
	Content string        `json:"content` //内容
	User    `json:"user"` //匿名结构体,继承
}

//SmsResMes
