package model

type User struct {
	//和tag一致
	UserName   string `json:"userName"`
	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserStatus int    `json:"userStatus"` //用户状态，在线离线
}
