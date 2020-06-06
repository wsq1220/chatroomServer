package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"code.mypro.com/my_chat_room/chat_room/server/proto"
	"github.com/astaxie/beego/logs"
)

type Client struct {
	conn   net.Conn
	userId int
	buf    [8192]byte
}

// receive
func (p *Client) readPackage() (msg proto.Message, err error) {
	n, err := p.conn.Read(p.buf[0:4])
	if n != 4 || err != nil {
		logs.Error("client read head data failed, err: %v", err)
		return
	}
	logs.Info("read head data: %v", p.buf[0:4])

	packLen := binary.BigEndian.Uint32(p.buf[0:4])
	logs.Debug("receive len: %v", packLen)

	n, err = p.conn.Read(p.buf[0:packLen])
	if n != int(packLen) {
		errMsg := fmt.Sprintf("read body data failed, expect: %v, actual: %v", int(packLen), n)
		err = errors.New(errMsg)
		logs.Error(errMsg)
		return
	}
	logs.Info("received body data: %v", string(p.buf[0:packLen]))

	if err = json.Unmarshal(p.buf[0:packLen], &msg); err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}

	return
}

// send
func (p *Client) writePackage(data []byte) (err error) {
	if data == nil {
		return
	}

	packLen := uint32(len(data))
	binary.BigEndian.PutUint32(p.buf[0:4], packLen)

	if _, err = p.conn.Write(p.buf[0:4]); err != nil {
		logs.Error("write head data failed, err: %v", err)
		return
	}
	logs.Info("write head data [%v] succ!", string(p.buf[0:4]))

	n, err := p.conn.Write(data)
	if err != nil {
		logs.Error("write data failed, err: %v", err)
		return
	}

	if n != int(packLen) {
		errMsg := fmt.Sprintf("write data not finished, now: %v/%v", n, int(packLen))
		err = errors.New(errMsg)
		logs.Error(errMsg)
		return
	}

	return
}

func (p *Client) Process() (err error) {
	for {
		var msg proto.Message
		msg, err = p.readPackage()
		if err != nil {
			logs.Error("read package failed whrn processing: %v", err)
			// TODO
			// clientMgr.DelClient(p.UserId)
			return
		}

		err = p.ProcessMsg(msg)
		if err != nil {
			logs.Error("process msg failed, err: %v, will continue", err)
			continue
		}
	}
}

// 处理请求
func (p *Client) ProcessMsg(msg proto.Message) (err error) {
	switch msg.Cmd {
	case proto.UserLoginCmd:
		err = p.login(msg)
	case proto.UserRegisterCmd:
		err = p.register(msg)
	case proto.UserSendMessageCmd:
		err = p.processUserSendMsg(msg)
	default:
		errMsg := fmt.Sprintf("the [%v] is not supported messgae!", msg.Cmd)
		err = errors.New(errMsg)
		logs.Error(errMsg)
		return
	}
	return
}

func (p *Client) login(msg proto.Message) (err error) {
	defer func() {
		p.loginResp(err)
	}()

	logs.Debug("enter login, got msg: %v", msg)
	var loginData proto.Login
	err = json.Unmarshal([]byte(msg.Data), &loginData)
	if err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}

	_, err = mgr.Login(loginData.Id, loginData.Password)
	if err != nil {
		logs.Error("login failed, err: %v", err)
		return
	}
	logs.Info("user[%v] login succ!", loginData.Id)

	// 通知其他用户此登录用户上线了
	p.notifyOtherUserOnline(loginData.Id)

	return
}

func (p *Client) notifyOtherUserOnline(userId int) {
	for id, client := range clientMgr.onlineUsers {
		if id == userId {
			continue
		}
		client.NotifyUserOnline(userId)
	}
}

func (p *Client) NotifyUserOnline(userId int) {
	var respMsg proto.Message
	respMsg.Cmd = proto.UserStatusNotifyCmd

	var notify proto.UserStatusNotify
	notify.UserId = userId
	notify.Status = proto.UserStatusOnline

	notifyData, err := json.Marshal(notify)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	respMsg.Data = string(notifyData)

	data, err := json.Marshal(respMsg)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	err = p.writePackage(data)
	if err != nil {
		logs.Error("notify other user you online failed, err: %v", err)
		return
	}
}

func (p *Client) loginResp(err error) {
	var respMsg proto.Message
	respMsg.Cmd = proto.UserLoginResCmd

	var loginResp proto.LoginResp
	loginResp.StatusCode = 200

	userMap := clientMgr.GetAllUsers()
	// 所有在线的用户
	for userId, _ := range userMap {
		loginResp.User = append(loginResp.User, userId)
	}

	if err != nil {
		loginResp.StatusCode = 500
		loginResp.Error = fmt.Sprintf("%v", err)
	}

	data, err := json.Marshal(loginResp)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	respMsg.Data = string(data)
	respData, err := json.Marshal(respMsg)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	err = p.writePackage(respData)
	if err != nil {
		logs.Error("write login resp data failed, err: %v", err)
		return
	}
}

func (p *Client) register(msg proto.Message) (err error) {
	var register proto.Register
	if err = json.Unmarshal([]byte(msg.Data), &register); err != nil {
		return
	}

	err = mgr.Register(&register.User)
	if err != nil {
		return
	}

	return
}

func (p *Client) processUserSendMsg(msg proto.Message) (err error) {
	var sendMsgReq proto.SendMsgReq
	err = json.Unmarshal([]byte(msg.Data), &sendMsgReq)
	if err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}

	users := clientMgr.GetAllUsers()
	for userId, client := range users {
		if userId == sendMsgReq.UserId {
			continue
		}
		client.SendMsgToUser(sendMsgReq.UserId, sendMsgReq.Data)
	}
	return
}

func (p *Client) SendMsgToUser(userId int, text string) {
	var respMsg proto.Message
	respMsg.Cmd = proto.UserRecvMessageCmd

	var recvMsg proto.UserRecvMsgReq
	recvMsg.UserId = userId
	recvMsg.Data = text

	recvMsgData, err := json.Marshal(recvMsg)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	respMsg.Data = string(recvMsgData)

	data, err := json.Marshal(respMsg)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	err = p.writePackage(data)
	if err != nil {
		logs.Error("send message failed, err: %v", err)
		return
	}
}
