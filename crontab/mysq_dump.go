package crontab

import (
	"bytes"
	"fmt"
	"github.com/lilihaooo/orange/settings"
	"github.com/lilihaooo/orange/utils/file"
	"github.com/robfig/cron/v3"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// 每天0点0分运行一次mysql备份
func mysqlDump() {
	c := cron.New()
	_, err := c.AddFunc("@every 1000s", func() {
		run()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
	select {}
}

// 读取配置文件参数,备份mysql数据库
func run() {
	mysqlConf := settings.Conf.MySQLConfig
	// 生成文件名字
	fileName := mysqlConf.DB + "-" + time.Now().Format("2006-01-02@15:04:05") + ".sql"
	pwd, _ := os.Getwd()
	// 生成在文件路径
	filePath := pwd + "/storage/database/" + time.Now().Format("2006-01-02") + "/"
	if err := file.MkDir(filePath); err != nil {
		logWrite.Error("定时备份数据库脚本失败,生成文件失败1：", err.Error())
		fmt.Println("定时备份数据库脚本失败,生成文件失败1：", err.Error())
	}
	realFile := filePath + fileName
	os.OpenFile(realFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	execDump(mysqlConf.Host, strconv.Itoa(mysqlConf.Port), mysqlConf.User, mysqlConf.Password, mysqlConf.DB, realFile)
}

// 执行备份命令
func execDump(dbHost, dbPort, dbUser, dbPassword, dbName, filePath string) {
	// docker mysql容器备份  mysql_master为mysql容器名称
	// commandDump := "docker exec -i mysql_master mysqldump --opt -h " + dbHost + " -P " + dbPort + " -u" + dbUser + " -p" + dbPassword + " " + dbName + " > " + filePath

	// mysql备份
	commandDump := "mysqldump --opt -h " + dbHost + " -P " + dbPort + " -u" + dbUser + " -p" + dbPassword + " " + dbName + " > " + filePath
	fmt.Println(commandDump)

	cmd := exec.Command("/bin/sh", "-c", commandDump)
	// 获取报错输出内容
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	// 如果报错 记录日志
	if err != nil {
		logWrite.Error("定时备份数据库脚本失败2：", err.Error())
		fmt.Println("定时备份数据库脚本失败2：", err.Error())
	} else {
		fmt.Println("定时备份数据库脚本成功:", out.String())
	}
	return
}
