package router

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"log"
	"net/http"
	"orange/handler/admin/baseHandler"
	"orange/help"
	"orange/middleware/cors"
	"orange/middleware/jwt"
	logToDb "orange/middleware/log"
	"orange/settings"
	"syscall"
)

type Router struct {
	r      *gin.Engine
	g      *gin.RouterGroup
	config string
}

func InitAdminRouter(config *settings.HttpConfig) *Router {
	mode := gin.ReleaseMode
	if config.GinMode == "debug" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Recovery())
	log := logrus.New()
	// 设置日志输出 生产环境生成日志文件
	if mode != "debug" {
		log.SetOutput(help.NewLogFileWriter("backend", "api"))
	}
	//todo log
	r.Use(ginlogrus.Logger(log))
	//静态文件怎么处理
	r.StaticFS("/storage", http.Dir("storage"))
	// 使用跨域中间件
	r.Use(cors.Cors())
	// -----------------------不需要权限的路由-----------------------//
	// 获取token
	r.POST("/login", baseHandler.Login)
	// 上传文件
	//r.POST("/uploadFile", handlers.UploadFile)
	// 开启使用jwt中间件
	r.Use(jwt.JWT())
	// 日志中间件 日志入库
	r.Use(logToDb.Log())

	// 获取oss文件上传临时授权token
	//r.GET("/getOssStsToken", handlers.GetOssStsToken)
	// 退出登陆
	r.DELETE("/logout", baseHandler.Logout)
	//// 获取用户信息(需要放在jwt下面通过token获取登陆用户信息)
	//r.GET("/adminInfo", baseHandler.AdminInfo)
	// -----------------------------------------------------------//

	// 使用casbin权限认证中间件
	//r.Use(casbin_auth.Auth())
	r.Routes()
	group := r.Group("v1")
	ip := "0.0.0.0"
	port := "8080"
	if config.Ip != "" {
		ip = config.Ip
	}
	if config.AdminPort != "" {
		port = config.AdminPort
	}
	s := fmt.Sprintf("%s:%s", ip, port)
	fmt.Println("process run at :", s)
	router := &Router{
		r,
		group,
		s,
	}

	// 用户
	//router.user()
	// 管理员
	router.admin()
	return router
}

// 管理员相关路由
func (router *Router) admin() {
	router.g.Use()
	{
		// 添加用户
		router.g.POST("/administrator", baseHandler.AdminAdd)
		// 更新管理员
		router.g.PUT("/administrator", baseHandler.AdminUpdate)
		// 管理员列表 127.0.0.1:8888/v1/administrator
		router.g.GET("/administrator", baseHandler.AdminList)
		// 管理员删除
		router.g.DELETE("/administrator", baseHandler.AdminDelete)
		// 管理员恢复
		router.g.PATCH("/administrator", baseHandler.AdminRecover)
		// 系统日志查询
		router.g.GET("/administrator/logs", baseHandler.LogList)
	}
}
func (router *Router) Start() {
	server := endless.NewServer(router.config, router.r)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}
	panic(server.ListenAndServe())
}
