package main

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
)

type LogConf struct {
	LogPath  string
	LogLevel string
}

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogger() (err error) {
	logConfig := make(map[string]interface{})
	// logConfig["filename"] = myConf.logConf.LogPath
	logConfig["filename"] = "./log/app.log"
	// logConfig["level"] = convertLogLevel(myConf.logConf.LogLevel)
	logConfig["level"] = convertLogLevel("debug")

	configData, err := json.Marshal(logConfig)
	if err != nil {
		logs.Error("json marshal failed, err: %v", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configData))
	// logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCall(true)

	return
}
