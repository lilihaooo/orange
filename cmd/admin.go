package cmd

import (
	"github.com/spf13/cobra"
	"orange/crontab"
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
		initCommonAction() // 配置文件, 数据库, redis
		//初始化后台工作
		initAdminAction()
		// 开启路由
		r := router.InitAdminRouter(settings.Conf.HttpConfig)
		r.Start()
	},
}

// 初始化后台工作
// 管理后台初始化操作
func initAdminAction() {
	// 定时任务开启 因为api和admin在同一个项目  所以定时任务只在后台开启
	crontab.InitCronJob()
}
