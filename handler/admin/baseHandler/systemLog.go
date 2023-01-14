package baseHandler

import (
	"github.com/gin-gonic/gin"
	"orange/handler"
	"orange/help"
	"orange/models/baseModel"
)

func LogList(c *gin.Context) {
	//分页map
	search := help.SearchParamsFormat(c)
	search["start_time"] = c.Query("start_time")
	search["end_time"] = c.Query("end_time")

	//model := baseModel.NewAdminApiLogModel()
	model := baseModel.AdminApiLog{}
	admin, total, err := model.GetLogList(search)
	if err != nil {
		handlers.FailWithSystem(c, err)
		return
	}
	response := make(map[string]interface{})
	response["list"] = admin
	response["total"] = total
	handlers.Success(c, response)

}

func LogStat(c *gin.Context) {
	model := baseModel.AdminApiLog{}
	data, err := model.StatLog()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, data)
}
