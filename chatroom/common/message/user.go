package message

//定义一个用户的结构体
type User struct {
	//为了序列化和反序列化成功
	//用户信息的json字符串必须和结构体字段的tag保持一致
	UserId     int    `json:userId`
	UserPwd    string `json:uesrPwd`
	UserName   string `json:userName`
	UserStatus int    `json:userStatus` //用户状态...
}
