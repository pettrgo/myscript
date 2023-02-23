package loggers

import (
	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	"myscript/config"
	"path/filepath"
)

func Init() {
	conf := config.GetConfig()
	serviceName := conf.ServiceName
	path := conf.Logger.Path
	level := conf.Logger.Level
	appFile := filepath.Join(path, serviceName+".app.log")
	errFile := filepath.Join(path, serviceName+".err.log")
	crashFile := filepath.Join(path, serviceName+".crash.log")
	err := logger.ResetStandardWithOptions(
		logger.Options{
			Level:     level,
			File:      appFile,
			ErrFile:   errFile,
			CrashFile: crashFile,
			AppName:   serviceName,
		})
	if err != nil {
		panic(err)
	}
}
