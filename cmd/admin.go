package cmd

import (
	"github.com/spf13/cobra"
	"orange/router"
	"orange/settings"
)

func init() {
	rootCmd.AddCommand(adminCmd)
}

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "orange后台api接口",
	Long:  `orange后台api接口`,
	Run: func(cmd *cobra.Command, args []string) {
		// 开启功能
		// 初始化必要工作
		initCommonAction() //配置文件, 数据库, redis
		// 开启路由
		r := router.InitAdminRouter(settings.Conf.HttpConfig)
		r.Start()
	},
}
