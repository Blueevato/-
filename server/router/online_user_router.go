package router

import "fmt"

// 因为OnlineUserRouter实例在服务器端有且只有一个
// 但很多地方都会用到，所以定义为全局变量
var (
	onlineUserRouter *OnlineUserRouter
)

// 在线用户 服务器端维护一个onlinemap
type OnlineUserRouter struct {
	onlineUsers map[int]*UserRouter
}

// 初始化
func init() {
	onlineUserRouter = &OnlineUserRouter{
		onlineUsers: make(map[int]*UserRouter, 1024),
	}
}

func (this *OnlineUserRouter) AddOnlineUser(ur *UserRouter) {
	this.onlineUsers[ur.UserId] = ur //添加
}

func (this *OnlineUserRouter) DelOnlineUser(userId int) {
	delete(this.onlineUsers, userId) //删除
}

func (this *OnlineUserRouter) GetAllOnlineUser() map[int]*UserRouter {
	return this.onlineUsers //返回全部在线用户
}
func (this *OnlineUserRouter) GetOnlineUserById(userId int) (ur *UserRouter, err error) {
	ur, ok := this.onlineUsers[userId]
	if !ok {
		//用户不在线
		err = fmt.Errorf("用户%d 不在线", userId)
		return
	}
	return //返回全部在线用户
}
