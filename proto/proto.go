package proto

// 这里会出现循环导入的问题
// import "github.com/wsq1220/chatroomServer/model"

type Message struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type Login struct {
	Id       int    `json:"user_id"`
	Password string `json:"password"`
}

type LoginResp struct {
	StatusCode int    `json:"status_code"`
	User       []int  `json:"users"`
	Error      string `json:"error"`
}

type Register struct {
	User User `json:"user"`
}

type SendMsgReq struct {
	UserId int    `json:"user_id"`
	Data   string `json:"data"`
}

type UserRecvMsgReq struct {
	UserId int    `json:"user_id"`
	Data   string `json:"data"`
}

type UserStatusNotify struct {
	UserId int `json:"user_id"`
	Status int `json:"user_status"`
}
