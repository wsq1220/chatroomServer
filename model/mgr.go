package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/wsq1220/chatroomServer/proto"
)

var (
	Table = "user"
)

// 其实是对redis pool的管理
// 初始化后可以使用
//  conn := pool.Get()
type UserMgr struct {
	pool *redis.Pool
}

func NewUserMgr(pool *redis.Pool) (userMgr *UserMgr) {
	userMgr = &UserMgr{
		pool: pool,
	}

	return
}

func (p *UserMgr) GetUser(conn redis.Conn, userId int) (user *proto.User, err error) {
	// defer conn.Close()
	res, err := redis.String(conn.Do("HGET", Table, fmt.Sprintf("%v", userId)))
	if err != nil {
		if err == redis.ErrNil {
			err = proto.ErrUserNotExist
		}
		// err = errGet
		logs.Error("test: get user failed, err: %v", err)
		return
	}
	logs.Info("result from getUser: %v", res)

	user = &proto.User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		logs.Error("json unmarshal failed, err: %v", err)
		return
	}
	logs.Info("got user: %v", user)

	return
}

func (p *UserMgr) Login(userId int, passwd string) (user *proto.User, err error) {
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

// register user to redis
func (p *UserMgr) Register(user *proto.User) (err error) {
	logs.Info("the param user: %v", user)
	conn := p.pool.Get()
	defer conn.Close()

	if user == nil {
		err = proto.ErrInvalidParam
		return
	}

	_, err = p.GetUser(conn, user.UserId)
	// logs.Error("got error when registering: %v", err)
	if err == nil {
		logs.Error("when err=nil, err: %v", err)
		err = proto.ErrUserExist
		return
	}

	if err != proto.ErrUserNotExist {
		logs.Error("test the error: %v", err)
		return
	}

	data, errJson := json.Marshal(user)
	if errJson != nil {
		err = errJson
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	_, err = conn.Do("HSet", Table, fmt.Sprintf("%v", user.UserId), string(data))
	if err != nil {
		logs.Error("set value to redis failed, err: %v", err)
		return
	}

	logs.Error("will exit, err: %v", err)
	return
}
