package main

import (
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

type Conf struct {
	redisConf RedisConf
	logConf   LogConf
}

var (
	// redisConf *RedisConf
	myConf *Conf
)

func loadLogConf(conf config.Configer) {
	myConf.logConf.LogPath = conf.String("log::log_path")
	myConf.logConf.LogLevel = conf.String("log::log_level")
	logs.Debug("log path: %v, log level: %v", myConf.logConf.LogPath, myConf.logConf.LogLevel)
}

// init redis conf and log conf
func initConf(confType, fileName string) (err error) {
	conf, err := config.NewConfig(confType, fileName)
	if err != nil {
		logs.Error("new config failed, err: %v", err)
		return
	}
	logs.Debug("new config succ!")

	// redisConf = &RedisConf{}
	myConf = &Conf{}
	myConf.redisConf.Addr = conf.String("redis::redis_addr")
	if len(myConf.redisConf.Addr) == 0 {
		logs.Warn("got redis addr failed,and will use default!")
		myConf.redisConf.Addr = "127.0.0.1:6379"
	}
	logs.Info("redis addr: %v", myConf.redisConf.Addr)

	myConf.redisConf.Password = conf.String("redis::redis_passwd")
	if len(myConf.redisConf.Password) == 0 {
		logs.Warn("not got redis password!")
	}
	logs.Info("got redis password: %v", myConf.redisConf.Password)

	myConf.redisConf.IdleConn, err = conf.Int("redis::redis_idle_conn")
	if err != nil {
		logs.Warn("not got redis idle conn:%v will use default", err)
		myConf.redisConf.IdleConn = 16
	}
	logs.Info("redis idle conn: %v", myConf.redisConf.IdleConn)

	myConf.redisConf.MaxConn, err = conf.Int("redis::redis_max_conn")
	if err != nil {
		logs.Warn("not got redis max conn:%v, will use default", err)
		myConf.redisConf.MaxConn = 1024
	}
	logs.Info("redis max conn: %v", myConf.redisConf.MaxConn)

	myConf.redisConf.IdleTimeout, err = conf.Int("redis::redis_time_out")
	if err != nil {
		logs.Warn("not got redis idle timeout: %v, will use default", err)
		myConf.redisConf.IdleConn = 300
	}
	logs.Info("redis idle timeout: %v", myConf.redisConf.IdleTimeout)

	// load log conf
	loadLogConf(conf)

	return
}
