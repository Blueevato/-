package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gochat2/common/model"
)

// 服务器启动时，就初始化一个UserDao实例
// 全局变量，可以直接操作
var (
	MyUserDao *UserDao
)

type UserDao struct {
	Pool *redis.Pool
}

// 工厂模式,创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		Pool: pool,
	}
	return
}

func (this *UserDao) getUserById(conn redis.Conn, id int) (user *model.User, err error) {
	//给定id去redis中查询用户
	//id 和userId绑定，为同一值
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			//users hash中无对应id
			err = model.ERROR_USER_NOTEXISTS
		}
		return
	}
	//反序列化成User
	err1 := json.Unmarshal([]byte(res), &user)
	if err1 != nil {
		fmt.Println("unmarshal err", err1)
		return
	}
	return

}

func (this *UserDao) registerUser(conn redis.Conn, userId int, userPwd, userName string) (err error) {
	//id 和userId绑定，为同一值
	// 将用户数据序列化为 JSON 存储
	userData := map[string]interface{}{
		"userId":   userId,
		"userPwd":  userPwd,
		"userName": userName,
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	_, err = conn.Do("HSET", "users", userId, jsonData)
	if err != nil {
		fmt.Println("register redis err", err)
		return
	}
	return
}

// 登录的校验
func (this *UserDao) Login(userId int, userPwd string) (user *model.User, err error) {

	//从redis连接池取出一个连接
	conn := this.Pool.Get()
	defer conn.Close()
	//用户是否存在
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	//校验用户的密码
	if user.UserPwd != userPwd {
		//失败
		err = model.ERROR_USER_PWD
		return
	}
	return
}

// 注册
func (this *UserDao) Register(user model.User) (err error) {

	//从redis连接池取出一个连接
	conn := this.Pool.Get()
	defer conn.Close()
	//用户是否存在
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = model.ERROR_USER_EXISTS
		return
	}

	//允许注册
	err = this.registerUser(conn, user.UserId, user.UserPwd, user.UserPwd)
	if err != nil {
		return
	}
	return
}
