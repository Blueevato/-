package main

import (
	"fmt"
	"gochat2/server/dao"
	"net"
)

// 处理和客户端的通讯
func process(conn net.Conn) {
	defer conn.Close()
	//调用总router
	pr := &ProcessRouter{
		Conn: conn,
	}
	err := pr.MainProcess()
	if err != nil {
		fmt.Println("main router err", err)
		return
	}
}

func initUserDao() {
	//这里要注意初始化的顺序，initPool -》 initUserDao
	dao.MyUserDao = dao.NewUserDao(pool) //这里的pools 是全局变量
}

func init() {
	//初始化连接池
	initPool("localhost:6379", 16, 0, 100)
	initUserDao()
}

func main() {

	//提示
	fmt.Println("服务器[new]在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("listen err")
	}

	for {
		fmt.Println("等待客户端连接")
		conn, err1 := listen.Accept()
		if err1 != nil {
			fmt.Println("accept err")
		}

		//链接成功，go一个协程和客户端保持通讯
		go process(conn)

	}
}
