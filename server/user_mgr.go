package main

import "code.mypro.com/my_chat_room/chat_room/server/model"

var (
	mgr *model.UserMgr
)

func initUserMgr() {
	mgr = model.NewUserMgr(pool)
}
