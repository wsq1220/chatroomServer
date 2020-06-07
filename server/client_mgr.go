package main

import (
	"fmt"

	"github.com/astaxie/beego/logs"
)

/*
	所有在线的用户
*/
type ClientMgr struct {
	onlineUsers map[int]*Client
}

var (
	clientMgr *ClientMgr
)

func init() {
	clientMgr = &ClientMgr{
		onlineUsers: make(map[int]*Client, 1024),
	}
}

func (p *ClientMgr) AddClient(userId int, client *Client) {
	p.onlineUsers[userId] = client
}

func (p *ClientMgr) GetClient(userId int) (client *Client, err error) {
	client, ok := p.onlineUsers[userId]
	if !ok {
		err = fmt.Errorf("user [%v] not exists!", userId)
		logs.Error(err.Error())
		return
	}

	return
}

// return map
func (p *ClientMgr) GetAllUsers() map[int]*Client {
	return p.onlineUsers
}

func (p *ClientMgr) DelClient(userId int) {
	delete(p.onlineUsers, userId)
}
