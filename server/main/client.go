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
