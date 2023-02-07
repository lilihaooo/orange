package couponModel

import "time"

type CouponIssuanceLog struct {
	ID        int64
	CouponId  int64      `json:"coupon_id" validate:"required"`  // 优惠券ID
	Num       int64      `json:"num" validate:"required"`        // 添加数量
	PreNum    int64      `json:"pre_num" validate:"required"`    // 添加前数量
	PostNum   int64      `json:"post_num" validate:"required"`   // 添加后数量
	AdminId   int64      `json:"admin_id" validate:"required"`   // adminId
	CreatedAt *time.Time `json:"created_at" validate:"required"` // 添加时间
}

// CouponList 优惠券列表
func (c *CouponIssuanceLog) CouIssLogList(params map[string]interface{}) (list []CouponIssuanceLog, count int, err error) {
	db := conn.Model(&CouponIssuanceLog{})
	db.Count(&count)
	err = db.Offset((params["search"].(int) - 1) * params["pageSize"].(int)).
		Limit(params["pageSize"]).
		Find(&list).Error
	return list, count, err

}
