package main

import (
	"fmt"
	"net"
	"redis5/chatroom/server/model"
	"time"
)

func precess(conn net.Conn) {
	//需要延迟关闭conn
	defer conn.Close()
	//循环的读客户端发送的信息
	processor := &Processor{
		Conn: conn,
	}
	err := processor.Fasongprocess()
	if err != nil {
		fmt.Printf("p.Fasong err: %v\n", err)
		return
	}
}

//编写一个函数完成对UserDao的初始化任务
func initUserDao() {
	//这里的pool本身就是全局变量
	//initPool必须在initUserDao前执行
	model.MyUserDao = model.NewUserDao(pool)
}
func initPoolAndUD() {
	//当服务器启动时我们就初始化链接池
	InitPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}
func main() {
	initPoolAndUD()
	//提示信息
	fmt.Println("服务器[新结构]在端口8889监听......")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Printf("net.listen err: %v\n", err)
		return
	}
	defer listen.Close()
	//一旦监听 就等待客户端连接服务端
	for {
		fmt.Println("等待客户端来链接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("listen.accept err: %v\n", err)
		}
		go precess(conn)
	}
}
