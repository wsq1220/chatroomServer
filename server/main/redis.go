package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

type RedisConf struct {
	Addr        string
	Password    string
	IdleConn    int
	MaxConn     int
	IdleTimeout int
	DbNum       int
}

func initRedis(addr string, idleConn, maxConn int, timeout time.Duration) (err error) {
	pool = &redis.Pool{
		MaxIdle:     idleConn,
		MaxActive:   maxConn,
		IdleTimeout: timeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}

	conn := pool.Get()
	defer conn.Close()

	// if has password
	// if _, err = conn.Do("AUTH", myConf.redisConf.Password); err != nil {
	// 	logs.Error("auth redis failed, err: %v", err)
	// 	return
	// }

	if _, err = conn.Do("ping"); err != nil {
		logs.Error("ping redis failed, err: %v", err)
		return
	}

	if _, err = conn.Do("SELECT", myConf.redisConf.DbNum); err != nil {
		logs.Error("select redis db num failed, err: %v", err)
		return
	}
	return
}
