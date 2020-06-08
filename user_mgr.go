package main

import "github.com/wsq1220/chatroomServer/model"

var (
	mgr *model.UserMgr
)

func initUserMgr() {
	mgr = model.NewUserMgr(pool)
}
