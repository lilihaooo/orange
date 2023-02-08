package search

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// 根据列表查询的参数,将其改成整形。page从0开始查询
func SearchParamsFormat(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})
	// 页码
	page := c.DefaultQuery("page", "1")
	params["page"], _ = strconv.Atoi(page)
	// 页数
	pageSize := c.DefaultQuery("pageSize", "10")
	params["pageSize"], _ = strconv.Atoi(pageSize)
	return params
}
