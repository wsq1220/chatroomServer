package model

import (
	"encoding/json"
	"fmt"
	"time"

	"code.mypro.com/my_chat_room/chat_room/server/proto"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

var (
	Table = "user"
)

// 其实是对redis pool的管理
type UserMgr struct {
	pool *redis.Pool
}

func NewUserMgr(pool *redis.Pool) (userMgr *UserMgr) {
	userMgr = &UserMgr{
		pool: pool,
	}

	return
}

func (p *UserMgr) GetUser(conn redis.Conn, userId int) (user *User, err error) {
	defer conn.Close()
	res, errGet := redis.String(conn.Do("HGET", Table, fmt.Sprintf("%v", userId)))
	if errGet != nil {
		if errGet == redis.ErrNil {
			err = proto.ErrUserNotExist
		}
		err = errGet
		logs.Error("get user failed, err: %v", err)
	}
	logs.Info("result from getUser: %v", res)

	user = &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}
	logs.Info("got user: %v", user)

	return
}

func (p *UserMgr) login(userId int, passwd string) (user *User, err error) {
	conn := p.pool.Get()
	defer conn.Close()

	user, err = p.GetUser(conn, userId)
	if err != nil {
		logs.Error("get user failed, err: %v", err)
		return
	}

	if user.UserId != userId || user.Password != passwd {
		err = proto.ErrInvalidPasswd
		return
	}

	user.Status = proto.UserStatusOnline
	user.LastLogin = fmt.Sprintf("%v", time.Now())

	return
}
