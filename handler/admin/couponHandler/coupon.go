package couponHandler

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/lilihaooo/orange/handler"
	"github.com/lilihaooo/orange/models/couponModel"
	search2 "github.com/lilihaooo/orange/utils/search"
)

type ICoupon interface {
	Add(c *gin.Context)
	Del(c *gin.Context)
	Edit(c *gin.Context)
	List(c *gin.Context)
	StateChange(c *gin.Context)
}

type Coupon struct {
}

func NewCoupon() ICoupon {
	return Coupon{}
}

func (c2 Coupon) Add(c *gin.Context) {
	coupon := new(couponModel.Coupon)
	err := c.BindJSON(coupon)
	if err != nil {
		handlers.FailWithParams(c, err)
		return
	}
	//获得当前admin的id
	userId, ok := c.Get("userId")
	if !ok {
		handlers.FailWithMessage(c, "获得用户Id失败")
		return
	}
	coupon.AdminId = userId.(int64)

	err = coupon.CouponCreate()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func (c2 Coupon) Del(c *gin.Context) {
	id := c.Query("id")
	coupon := new(couponModel.Coupon)
	err := coupon.CouponDelete(id)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func (c2 Coupon) Edit(c *gin.Context) {
	coupon := new(couponModel.Coupon)
	err := c.BindJSON(coupon)
	if err != nil {
		handlers.FailWithParams(c, err)
		return
	}
	//获得当前admin的id
	userId, ok := c.Get("userId")
	if !ok {
		handlers.FailWithMessage(c, "获得用户Id失败")
		return
	}
	coupon.AdminId = userId.(int64)
	err = coupon.CouponUpdate()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func (c2 Coupon) List(c *gin.Context) {
	search := search2.SearchParamsFormat(c)
	// 查询是否有username、name、mobile等查询参数
	search["name"] = c.Query("name")
	search["admin_id"] = c.Query("admin_id")
	search["start_time"] = c.Query("start_time")
	search["end_time"] = c.Query("end_time")
	search["status"] = c.Query("status")
	search["sort_type"] = c.Query("sort_type")
	coupon := &couponModel.Coupon{}
	list, total, err := coupon.CouponList(search)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	response := make(map[string]interface{})
	response["list"] = list
	response["total"] = total
	handlers.Success(c, response)
}

func (c2 Coupon) StateChange(c *gin.Context) {
	id := c.Query("id")
	coupon := new(couponModel.Coupon)
	err := coupon.CouponStateChange(id)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}
