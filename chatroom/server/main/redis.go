package main

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

//定义全局的pool
var pool *redis.Pool

//当启动程序时就初始化连接池
func InitPool(address string, maxIdle int, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大空闲链接数
		MaxActive:   maxActive,   //表示已经和数据库的最大链接数，0表示没有限制
		IdleTimeout: idleTimeout, //最大空闲时间
		Dial: func() (redis.Conn, error) { //初始化链接代码，链接哪个IP
			return redis.Dial("tcp", address)
		},
	}
}
