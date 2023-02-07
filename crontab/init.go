package crontab

import (
	"fmt"
	"github.com/lilihaooo/orange/settings"
	log2 "github.com/lilihaooo/orange/utils/log"
	log "github.com/sirupsen/logrus"
)

var logWrite *log.Logger

func init() {
	log := log.New()
	// 设置日志输出
	log.SetOutput(log2.NewLogFileWriter("crontab", "error"))
	logWrite = log
}

func InitCronJob() {
	// 定时脚本任务,根据配置文件,只有一台服务器跑,如果配置返回false 不执行
	if !settings.Conf.EnvConfig.Crontab {
		fmt.Println("副服务器")
		return
	}
	fmt.Println("主服务器")
	// 备份数据库文件
	go mysqlDump()
	go couponIssuance()
}
