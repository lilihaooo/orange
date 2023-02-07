package cmd

import (
	"fmt"
	"github.com/lilihaooo/orange/db/conn/mysql"
	"github.com/lilihaooo/orange/db/conn/redis"
	"github.com/lilihaooo/orange/models/conn"
	"github.com/lilihaooo/orange/settings"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "orange",
	Short: "orange",
	Long:  `orange`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("more detail please input -h")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 应用配置信息、数据库、脚本初始化
func initCommonAction() {
	// 加载配置文件
	settings.InitConfig("/config/config.yaml")

	// 初始化数据库信息
	mysql.Connect(settings.Conf.MySQLConfig)
	conn.InitConn()

	//初始化redis连接
	redis.InitRedisConn()
}
