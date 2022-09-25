package model

import (
	"net"
	"redis5/chatroom/common/message"
)
//因为在客户端，很多地方会使用到curUser,我们将其作为一个全局

type CurUser struct {
	Conn net.Conn
	message.User
}
