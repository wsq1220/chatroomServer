package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
)

func main() {

	if err := initLogger(); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		panic(err)
	}
	fmt.Println("init logger succ!")

	err := initConf("ini", "./server/conf/app.conf")
	if err != nil {
		logs.Error("init config failed, err: %v", err)
		panic(err)
	}
	logs.Info("init conf succ!")

	err = initRedis(myConf.redisConf.Addr, myConf.redisConf.IdleConn, myConf.redisConf.MaxConn, time.Duration(myConf.redisConf.IdleTimeout)*time.Second)
	if err != nil {
		logs.Error("init redis failed, err: %v", err)
		panic(err)
	}
	logs.Info("init redis succ!")

	initUserMgr()

	runServer("0.0.0.0:10000")
}
