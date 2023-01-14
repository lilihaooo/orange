package crontab

import (
	"fmt"
	"orange/settings"
)

func InitCronJob() {
	// 定时脚本任务,根据配置文件,只有一台服务器跑,如果配置返回false 不执行
	if !settings.Conf.EnvConfig.Crontab {
		fmt.Println("副服务器")
		return
	}
	fmt.Println("主服务器")
	// 备份数据库文件
	go mysqlDump()
}
