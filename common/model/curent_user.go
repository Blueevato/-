package model

import (
	"net"
)

// 全局,当前用户
type CurentUser struct {
	Conn net.Conn
	User
}
