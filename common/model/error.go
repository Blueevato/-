package model

import "errors"

// 自定义错误
var (
	ERROR_USER_NOTEXISTS = errors.New("该用户不存在...")
	ERROR_USER_EXISTS    = errors.New("用户已注册...")
	ERROR_USER_PWD       = errors.New("用户密码错误...")
)
