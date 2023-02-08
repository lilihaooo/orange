package couponModel

import (
	"time"
)

type CouponIssuanceLog struct {
	ID        int64
	CouponId  int64      `json:"coupon_id" validate:"required" label:"优惠券ID"` // 优惠券ID
	Num       int64      `json:"num" validate:"required" label:"添加数量"`        // 添加数量
	PreNum    int64      `json:"pre_num" validate:"required" label:"添加前数量"`   // 添加前数量
	PostNum   int64      `json:"post_num" validate:"required" label:"添加后数量"`  // 添加后数量
	AdminId   int64      `json:"admin_id" validate:"required" label:"管理员ID"`  // adminId
	CreatedAt *time.Time `json:"created_at" validate:"required" label:"添加时间"` // 添加时间
}

// CouponList 优惠券列表
func (c *CouponIssuanceLog) CouIssLogList(params map[string]interface{}) (list []CouponIssuanceLog, count int64, err error) {
	db := conn.Model(&CouponIssuanceLog{})
	db.Count(&count)
	if params["begin_time"] != nil {
		db = db.Where("created_at >= ?", params["begin_time"])
	}
	if params["end_time"] != nil {
		db = db.Where("created_at <= ?", params["end_time"])
	}
	err = db.Offset((params["page"].(int) - 1) * params["pageSize"].(int)).
		Limit(params["pageSize"].(int)).
		Find(&list).Error
	return list, count, err
}
