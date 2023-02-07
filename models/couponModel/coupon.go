package couponModel

import (
	"errors"
	"github.com/lilihaooo/orange/models"
	string2 "github.com/lilihaooo/orange/utils/str"
	"github.com/lilihaooo/orange/utils/validCheck"
	"github.com/sirupsen/logrus"
	"strings"
)

type Coupon struct {
	models.Model
	Name        string  `json:"name" validate:"required,max=20"`      // 名称
	Description string  `json:"description" validate:"max=100"`       // 描述
	Money       float64 `json:"money" validate:"required,lte=10000"`  // 券面金额
	LowerMoney  float64 `json:"lower_money" validate:"lte=10000"`     // 启用金额
	AdminId     int64   `json:"admin_id" validate:"required"`         // 管理员id
	Tip         string  `json:"tip" validate:"max=50"`                // 券面提示
	Total       int64   `json:"total" validate:"required"`            // 已发数量
	UsedNum     int64   `json:"used_num" validate:"required"`         // 使用数量
	StartTime   string  `json:"start_time" validate:"required,max=5"` // 抢券开始时间
	EndTime     string  `json:"end_time" validate:"required,max=5"`   // 抢券结束时间
	Status      int     `json:"status" validate:"required,lte=2"`     // 状态
}

// CouponCreate 添加优惠券
func (c *Coupon) CouponCreate() (err error) {
	//验证数据
	err = validCheck.Validate(c)
	if err != nil {
		return err
	}
	return conn.Create(&c).Error
}

// CouponDelete 删除优惠券
func (c *Coupon) CouponDelete(id string) (err error) {
	return conn.Where("id = ?", id).Delete(&c).Error
}

// CouponUpdate 修改优惠券信息
func (c *Coupon) CouponUpdate() (err error) {
	//验证数据
	err = validCheck.Validate(c)
	if err != nil {
		return err
	}
	newCoupon := new(Coupon)
	if conn.Where("id = ?", c.ID).Find(newCoupon).RowsAffected == 0 {
		return errors.New("优惠券不存在")
	}

	newCoupon.ID = c.ID
	newCoupon.Description = c.Description
	newCoupon.Money = c.Money
	newCoupon.LowerMoney = c.LowerMoney
	newCoupon.AdminId = c.AdminId
	newCoupon.Tip = c.Tip
	newCoupon.Total = c.Total
	newCoupon.StartTime = c.StartTime
	newCoupon.EndTime = c.EndTime
	newCoupon.Status = c.Status
	return conn.Save(newCoupon).Error
}

// CouponList 优惠券列表
func (c *Coupon) CouponList(params map[string]interface{}) (list []Coupon, count int, err error) {
	db := conn.Model(&Coupon{})
	// 筛选
	if params["name"] != "" {
		db = db.Where("name LIKE ?", "%"+params["name"].(string)+"%")
	}
	if params["admin_id"] != "" {
		db = db.Where("admin_id = ? ", params["admin_id"])
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
	   1. 券面金额降序
	   2. 券面金额升序
	   3. 启用金额降序
	   4. 启用金额升序
	   5. 已发数量降序
	   6. 已发数量升序
	   7. 使用数量降序
	   8. 使用数量升序
	   7. id降序
	   8. id升序
	*/
	var sortStr string
	if params["sort_type"] != "0" {
		sortStr = "id ASC"
	}
	if params["sort_type"] == "1" {
		sortStr = "money DESC"
	}
	if params["sort_type"] == "2" {
		sortStr = "money ASC"
	}
	if params["sort_type"] == "3" {
		sortStr = "lower_money DESC"
	}
	if params["sort_type"] == "4" {
		sortStr = "lower_money ASC"
	}
	if params["sort_type"] == "5" {
		sortStr = "total DESC"
	}
	if params["sort_type"] == "6" {
		sortStr = "total ASC"
	}
	if params["sort_type"] == "7" {
		sortStr = "used_num DESC"
	}
	if params["sort_type"] == "8" {
		sortStr = "used_num ASC"
	}
	if params["sort_type"] == "9" {
		sortStr = "id DESC"
	}
	if params["sort_type"] == "10" {
		sortStr = "id ASC"
	}

	db.Count(&count)
	err = db.Offset((params["search"].(int) - 1) * params["pageSize"].(int)).
		Limit(params["pageSize"]).
		Order(sortStr).
		Find(&list).Error
	return list, count, err

}

// CouponStateChange 优惠券状态修改
func (c *Coupon) CouponStateChange(id string) (err error) {
	if conn.Where("id = ?", id).Find(&c).RowsAffected == 0 {
		return errors.New("优惠券不存在")
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

type Issue struct {
	CouponId  int64
	TimePoint string
	Num       int64
}

// IssueCoupons 发放优惠券
func IssueCoupons(logWrite *logrus.Logger) error {
	// 获得当前年月日
	date := string2.CurrentTimeYMD()
	// 获得时分
	time := string2.CurrentTimeHI()
	time = "17:30"
	// 开始时间结束时间包含当前日期的预发券ids
	var issue []Issue
	sql := `select 
coupon_pre_issuance.coupon_id, time_point, num
from coupon_pre_issuance 
inner join 
coupon 
on coupon.id = coupon_pre_issuance.coupon_id 
where coupon_pre_issuance.start_time <= ? and coupon_pre_issuance.end_time > ? 
and coupon_pre_issuance.status = 1 
and coupon.status = 1`
	err := conn.Raw(sql, date, date).Scan(&issue).Error
	if err != nil {
		return err
	}
	for _, item := range issue {
		//将时间点字符串截断为每个时间点的切片
		timePointSlice := strings.Split(item.TimePoint, ",")
		for _, timePoint := range timePointSlice {
			if timePoint == time {
				err = action(&item)
				if err != nil {
					//todo 记录日志
					logWrite.Error("发券脚本执行失败:" + err.Error())
				}
			}
		}
	}
	return nil
}

func action(issue *Issue) error {

	coupon := Coupon{}
	conn.Where("id = ?", issue.CouponId).First(&coupon)
	// 获得该券发放前的数量
	preNum := coupon.Total
	// 发券后数量
	coupon.Total = coupon.Total + issue.Num

	// 开启事务
	tx := conn.Begin()
	// 1. 修改coupon表发放数量
	err := tx.Save(&coupon).Error
	if err != nil {
		return err
	}
	// 2. 新增日志
	logInfo := CouponIssuanceLog{}
	logInfo.CouponId = issue.CouponId
	logInfo.Num = issue.Num
	logInfo.AdminId = 1 // todo 应该是添加预发信息的admin, 但是预发表没有这个字段!!
	logInfo.PreNum = preNum
	logInfo.PostNum = coupon.Total
	err = tx.Create(&logInfo).Error
	//err = errors.New("测试事务")
	if err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	tx.Commit()
	return nil
}
