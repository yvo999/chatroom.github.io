package model

import "errors"

//根据业务逻辑需要自定义错误
var (
	ERROR_USER_NOTEXISTS = errors.New("用户不存在...")
	ERROR_USER_EXISTS    = errors.New("用户已存在...")
	ERROR_USER_PWD       = errors.New("密码不正确...")
)
