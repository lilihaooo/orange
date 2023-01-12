package help

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"time"
)

// RandString 生成随机字符串
func RandString(len int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// 获取年月日
func CurrentTimeYMD() string {
	return time.Now().Format("2006-01-02")
}

// 根据列表查询的参数,将其改成整形。page从0开始查询
func SearchParamsFormat(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})
	// 页码
	page := c.DefaultQuery("page", "1")
	params["page"], _ = strconv.Atoi(page)
	params["page"] = params["page"]
	// 页数
	pageSize := c.DefaultQuery("pageSize", "10")
	params["pageSize"], _ = strconv.Atoi(pageSize)
	return params
}

// 获取年月日时分秒
func CurrentTimeYMDHIS() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
