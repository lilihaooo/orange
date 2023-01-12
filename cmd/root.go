package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"orange/db/conn/mysql"
	"orange/db/conn/redis"
	"orange/models/conn"
	"orange/settings"
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
