package processes

import (
	"fmt"
	"redis5/chatroom/client/model"
	"redis5/chatroom/common/message"
)

//客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)

//因为在客户端，很多地方会使用到curUser,我们将其作为一个全局
var CurUser model.CurUser //在用户登陆成功后，完成对其初始化

//在客户端显示当前在线的用户
func outputOnlineUser() {
	fmt.Println("当前在线用户列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户id: \t", id)
	}
}

//编写一个方法处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	//优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok { //原来没有
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	outputOnlineUser()
}
