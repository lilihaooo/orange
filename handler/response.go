package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log2 "github.com/lilihaooo/orange/utils/log"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = struct{}{}
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "success",
		"data":    data,
	})
	return
}

func FailWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    40000,
		"message": message,
	})
	return
}

func FailWithParams(c *gin.Context, error error) {
	fmt.Println(error)
	c.JSON(http.StatusOK, gin.H{
		"code":    40003,
		"message": "请求参数错误",
	})
	return
}

func FailWithSystem(c *gin.Context, err error) {
	fmt.Println(err)
	//log := log.New()
	// 设置日志输出
	log.SetOutput(log2.NewLogFileWriter("system", "error"))
	log.Error(err)
	c.JSON(http.StatusOK, gin.H{
		"code":    50000,
		"message": "系统内部错误",
	})
	return
}
