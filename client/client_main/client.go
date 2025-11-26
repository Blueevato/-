package main

import (
	"client/client_router"
	"fmt"
)

var (
	userId   int
	userPwd  string
	userName string
)

func main() {
	key := 0
	loop := false
	for !loop {
		fmt.Println("--menu--")
		fmt.Println("--\t\t\t 1 登录聊天室--")
		fmt.Println("--\t\t\t 2 注册--")
		fmt.Println("--\t\t\t 3 退出--")
		fmt.Println("--\t\t\t 请选择(1-3):--")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("--login--")
			case1()
		case 2:
			fmt.Println("--register--")
			case2()
		case 3:
			fmt.Println("--exit--")
			loop = true
		default:
			fmt.Println("--input err--")
		}
	}
}

func case1() {
	fmt.Println("input uid")
	fmt.Scanf("%d\n", &userId) //Scanln更优雅
	fmt.Println("input upwd")
	fmt.Scanf("%s\n", &userPwd)
	up := &client_router.ClientUserRouter{}
	err := up.Login(userId, userPwd)
	if err != nil {
		fmt.Println("login err")
	}
}

func case2() {
	fmt.Println("input uid")
	fmt.Scanf("%d\n", &userId) //Scanln更优雅
	fmt.Println("input upwd")
	fmt.Scanf("%s\n", &userPwd)
	fmt.Println("input uname")
	fmt.Scanf("%s\n", &userName)
	up := &client_router.ClientUserRouter{}
	err := up.Register(userId, userPwd, userName)
	if err != nil {
		fmt.Println("register err")
	}
}
