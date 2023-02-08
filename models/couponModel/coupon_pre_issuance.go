package couponModel

import (
	"errors"
	"github.com/lilihaooo/orange/models"
	"github.com/lilihaooo/orange/utils/validCheck"
	"strings"
)

//优惠券预发

type CouponPreIssuance struct {
	models.Model
	CouponId  int64  `json:"coupon_id" validate:"required" label:"优惠券ID"`          // 优惠券Id
	StartTime string `json:"start_time" validate:"required,max=20" label:"发放开始时间"` // 发放开始时间
	EndTime   string `json:"end_time" validate:"required,max=20" label:"发放结束时间"`   // 发放结束时间
	Status    int    `json:"status" validate:"required,max=2" label:"状态"`          // 状态
	TimePoint string `json:"time_point" validate:"required,max=255" label:"发放时间点"` // 发放时间点
	Num       int64  `json:"num" validate:"required,min=0" label:"单次发放数量"`         // 单次发放数量
	Coupon    Coupon `gorm:"FOREIGNKEY:coupon_id;ASSOCIATION_FOREIGNKEY:ID" validate:"-"`
}

// CouPreIssCreate 添加优惠券预发
func (c *CouponPreIssuance) CouPreIssCreate() (err error) {
	//验证数据
	err = validCheck.Validate(c)
	if err != nil {
		return err
	}
	if !validCheck.StrISYmd(c.StartTime) || !validCheck.StrISYmd(c.EndTime) {
		return errors.New("日期格式不正确")
	}
	if conn.Where("id = ?", c.CouponId).Find(&[]Coupon{}).RowsAffected == 0 {
		return errors.New("优惠券ID不存在")
	}

	// 验证发券时间点
	timePoint := strings.Split(c.TimePoint, ",")
	for _, item := range timePoint {
		if !validCheck.StrISMs(item) {
			return errors.New("发放时间点格式错误")
		}
	}
	// 验证发券开始时间和结束时间

	return conn.Create(&c).Error
}

// CouPreIssDelete 删除优惠券预发
func (c *CouponPreIssuance) CouPreIssDelete(id string) (err error) {
	return conn.Where("id = ?", id).Delete(&c).Error
}

// CouPreIssUpdate 修改优惠券信息预发
func (c *CouponPreIssuance) CouPreIssUpdate() (err error) {
	err = validCheck.Validate(c)
	if err != nil {
		return err
	}
	newCouponPreIssuance := new(CouponPreIssuance)
	if conn.Where("id = ?", c.ID).First(newCouponPreIssuance).RowsAffected == 0 {
		return errors.New("记录不存在")
	}
	newCouponPreIssuance.ID = c.ID
	newCouponPreIssuance.TimePoint = c.TimePoint
	newCouponPreIssuance.StartTime = c.StartTime
	newCouponPreIssuance.EndTime = c.EndTime
	newCouponPreIssuance.Num = c.Num
	newCouponPreIssuance.Status = c.Status
	return conn.Save(newCouponPreIssuance).Error
}

// CouPreIssList 优惠券列表预发
func (c *CouponPreIssuance) CouPreIssList(params map[string]interface{}) (list []CouponPreIssuance, count int64, err error) {
	db := conn.Model(&CouponPreIssuance{})
	// 筛选
	if params["coupon_id"] != "" {
		db = db.Where("coupon_id = ? ", params["coupon_id"])
	}
	if params["start_time"] != "" {
		db = db.Where("created_at >= ? ", params["start_time"])
	}
	if params["end_time"] != "" {
		db = db.Where("created_at <= ? ", params["end_time"])
	}
	if params["status"] != "" {
		db = db.Where("status = ? ", params["status"])
	}

	// 排序
	/*
	   SortType排序规则:
	   1. 单次发放数量降序
	   2. 单次发放数量升序
	   3. Id降序
	   4. Id升序
	   5. 开始时间降序
	   6. 开始时间升序
	   7. 结束时间降序
	   8. 结束时间升序

	*/
	var sortStr string
	if params["sort_type"] != "0" {
		sortStr = "id ASC"
	}
	if params["sort_type"] == "1" {
		sortStr = "num DESC"
	}
	if params["sort_type"] == "2" {
		sortStr = "num ASC"
	}
	if params["sort_type"] == "3" {
		sortStr = "id DESC"
	}
	if params["sort_type"] == "4" {
		sortStr = "id ASC"
	}
	if params["sort_type"] == "5" {
		sortStr = "begin_time DESC"
	}
	if params["sort_type"] == "6" {
		sortStr = "begin_time ASC"
	}
	if params["sort_type"] == "7" {
		sortStr = "end_time DESC"
	}
	if params["sort_type"] == "8" {
		sortStr = "end_time ASC"
	}

	db.Count(&count)
	err = db.Offset((params["page"].(int) - 1) * params["pageSize"].(int)).Preload("Coupon").
		Limit(params["pageSize"].(int)).
		Order(sortStr).
		Find(&list).Error
	return list, count, err

}

// CouPreIssStateChange 优惠券预发状态修改
func (c *CouponPreIssuance) CouPreIssStateChange(id string) (err error) {
	if conn.Where("id = ?", id).First(&c).RowsAffected == 0 {
		return errors.New("记录不存在")
	}
	var now int
	if c.Status == 1 {
		now = 2
	}
	if c.Status == 2 {
		now = 1
	}
	return conn.Model(&c).Where("id = ?", id).Update("status", now).Error
}
