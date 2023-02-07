package couponHandler

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/lilihaooo/orange/handler"
	"github.com/lilihaooo/orange/models/couponModel"
	search2 "github.com/lilihaooo/orange/utils/search"
)

func CouIssLogList(c *gin.Context) {
	search := search2.SearchParamsFormat(c)
	model := &couponModel.CouponIssuanceLog{}
	list, total, err := model.CouIssLogList(search)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	response := make(map[string]interface{})
	response["list"] = list
	response["total"] = total
	handlers.Success(c, response)

}
