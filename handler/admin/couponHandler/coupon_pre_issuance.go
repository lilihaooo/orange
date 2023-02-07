package couponHandler

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/lilihaooo/orange/handler"
	"github.com/lilihaooo/orange/models/couponModel"
	search2 "github.com/lilihaooo/orange/utils/search"
)

type ICouponPreIssuance interface {
	Add(c *gin.Context)
	Del(c *gin.Context)
	Edit(c *gin.Context)
	List(c *gin.Context)
	StateChange(c *gin.Context)
}

type CouponPreIssuance struct {
}

func NewCouponPreIssuance() ICoupon {
	return CouponPreIssuance{}
}

func (c2 CouponPreIssuance) Add(c *gin.Context) {
	couponPreIssuance := new(couponModel.CouponPreIssuance)
	err := c.BindJSON(couponPreIssuance)
	if err != nil {
		handlers.FailWithParams(c, err)
		return
	}
	err = couponPreIssuance.CouPreIssCreate()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func (c2 CouponPreIssuance) Del(c *gin.Context) {
	id := c.Query("id")
	couponPreIssuance := new(couponModel.CouponPreIssuance)
	err := couponPreIssuance.CouPreIssDelete(id)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func (c2 CouponPreIssuance) Edit(c *gin.Context) {
	couponPreIssuance := new(couponModel.CouponPreIssuance)
	err := c.BindJSON(couponPreIssuance)
	if err != nil {
		handlers.FailWithParams(c, err)
		return
	}
	err = couponPreIssuance.CouPreIssUpdate()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func (c2 CouponPreIssuance) List(c *gin.Context) {
	search := search2.SearchParamsFormat(c)
	// 查询是否有username、name、mobile等查询参数
	search["coupon_id"] = c.Query("coupon_id")
	search["start_time"] = c.Query("start_time")
	search["end_time"] = c.Query("end_time")
	search["status"] = c.Query("status")
	search["sort_type"] = c.Query("sort_type")
	couponPreIssuance := &couponModel.CouponPreIssuance{}
	list, total, err := couponPreIssuance.CouPreIssList(search)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	response := make(map[string]interface{})
	response["list"] = list
	response["total"] = total
	handlers.Success(c, response)
}

func (c2 CouponPreIssuance) StateChange(c *gin.Context) {
	id := c.Query("id")
	couponPreIssuance := new(couponModel.CouponPreIssuance)
	err := couponPreIssuance.CouPreIssStateChange(id)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}
