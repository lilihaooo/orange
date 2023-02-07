package baseHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/lilihaooo/orange/handler"
	"github.com/lilihaooo/orange/models/baseModel"
	search2 "github.com/lilihaooo/orange/utils/search"
)

func LogList(c *gin.Context) {
	//分页map
	search := search2.SearchParamsFormat(c)
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

func LogStatistics(c *gin.Context) {
	model := baseModel.AdminApiLog{}
	data, err := model.StatisticsLog()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, data)
}
