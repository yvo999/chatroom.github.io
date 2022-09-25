package model

import (
	"encoding/json"
	"fmt"
	"redis5/chatroom/common/message"

	"github.com/gomodule/redigo/redis"
)

//我们在服务器启动后就初始化UserDao实例
//把它做成全局变量，在需要于redis操作时，就直接使用
var (
	MyUserDao *UserDao
)

//定义一个结构体完成对User结构体的操作
type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式获取一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//根据用户ID返回一个User实例+error
func (this *UserDao) GetUserById(conn redis.Conn, id int) (user *message.User, err error) {
	//通过给定的ID去redis查询用户
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		if err == redis.ErrNil { //表示在哈希表中没有找到对应ID
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	//把res反序列化成User实例
	user = &message.User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Printf("user j.S err: %v\n", err)
		return
	}
	return
}

//完成登陆效验
//1.Login 完成对用户的验证
//2.如果用户id&pwd都正确，则返回user实例
//3.如果有错误则返回对应错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *message.User, err error) {
	//先从UserDao的链接池取出一根链接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.GetUserById(conn, userId)
	if err != nil {
		return
	}
	//这是用户获取到了
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {
	//先从UserDao的链接池取出一根链接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.GetUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS //这是用户名被占用了
		return
	}
	//这时说明id在redis还没有，则完成注册
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("序列化入库err: %v\n", err)
		return
	}
	//入库
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误", err)
		return
	}
	return
}
