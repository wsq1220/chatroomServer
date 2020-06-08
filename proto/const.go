package proto

var (
	UserStatusOnline  = 1
	UserStatusOffline = 0
)

const (
	UserLoginCmd        = "user_login"
	UserLoginResCmd     = "user_Login_res"
	UserRegisterCmd     = "user_register"
	UserStatusNotifyCmd = "user_status_notify"
	UserSendMessageCmd  = "user_send_message"
	UserRecvMessageCmd  = "user_recv_message"
)
