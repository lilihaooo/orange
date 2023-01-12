package main

import (
	log "github.com/sirupsen/logrus"
	"orange/cmd" //执行保中的init
	"os"
)

func init() {
	//设置日志为Text格式
	log.SetFormatter(&log.TextFormatter{
		//禁用颜色
		//DisableColors: true,
		//完整时间戳
		FullTimestamp: true,
	})

	log.SetReportCaller(true)

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	log.SetFormatter(&log.TextFormatter{})

	// 日志级别debug
	log.SetLevel(log.TraceLevel)
}

func main() {
	cmd.Execute()
}
