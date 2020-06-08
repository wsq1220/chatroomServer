package main

import (
	"fmt"
	"net"

	"github.com/astaxie/beego/logs"
	"github.com/fatih/color"
)

func runServer(addr string) (err error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("net listen failed, err: %v", err)
		return
	}
	theNotice := fmt.Sprintf("server is listening: %v", addr)
	color.Set(color.FgGreen, color.Bold)
	// fmt.Printf("server is listening: %v\n", addr)
	fmt.Println(theNotice)
	color.Unset()

	for {
		conn, err := l.Accept()
		if err != nil {
			logs.Error("accept failed, err: %v, will continue", err)
			continue
		}

		go process(conn)
	}
}

func process(conn net.Conn) {
	// 捕获goroutine中的err
	defer func() {
		// 关闭句柄
		conn.Close()
		// 捕获错误
		if err := recover(); err != nil {
			logs.Error("error occoured in process goroutine: %v", err)
		}
	}()

	client := Client{
		conn: conn,
	}

	err := client.Process()
	if err != nil {
		logs.Error("client process failed, err: %v", err)
		return
	}
}
