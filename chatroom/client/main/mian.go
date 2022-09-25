package main

import (
	"fmt"
	"os"
	processes "redis5/chatroom/client/process"
)

//定义两个全局变量一个表示账号一个表示密码
var userId int
var userPwd string
var userName string

func main() {
	//接受用户的选择
	var key int
	//判断是否还继续菜单

	for {
		fmt.Println("----------欢迎登陆多人聊天系统----------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3)")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			//说明用户要登陆
			fmt.Println("请输入账号")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("输入用户密码")
			fmt.Scanf("%s\n", &userPwd)
			//先把登陆的函数，写到另一个文件里,如login.go
			up := &processes.UserPorces{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("输入昵称:")
			fmt.Scanf("%s\n", &userName)
			//调用Userproces,完成注册
			up := &processes.UserPorces{}
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)

		default:
			fmt.Println("你的输入有误请重新输入")
		}
	}
}
