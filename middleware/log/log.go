package log

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/lilihaooo/orange/models/baseModel"
	"io/ioutil"
	"time"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 查询操作不入库
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		// 字节流 转为string后就是前端输出的内容
		// 获取原body数据, body的内容只能读取一次，后面在读取都是空的。所以需要重新赋值
		body, _ := c.GetRawData()
		// 将原body塞回去  把读过的字节流重新放到body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		//fmt.Println(str(body))

		go func(c *gin.Context, content string) {
			// pprof 性能分析
			// 1.找文件  如果生成文件
			// 2.写入文件
			// 4.关闭文件
			//roles, _ := c.Get("userRole")
			userId, _ := c.Get("userId")
			Username, _ := c.Get("username")
			t := time.Now()
			//roleString := ""
			//for _, v := range roles.([]str) {
			//	roleString += v + ","
			//}
			// 添加到日志表
			tmp := baseModel.AdminApiLog{
				UserId:   userId.(int64),
				Username: Username.(string),
				//Role:      roleString,
				Host:      c.ClientIP(),
				Path:      c.Request.RequestURI,
				Method:    c.Request.Method,
				Content:   content,
				CreatedAt: &t,
			}
			tmp.AddAdminApiLog()
		}(c, string(body))

		c.Next()
	}
}
