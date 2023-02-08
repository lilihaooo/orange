package crontab

import (
	"github.com/lilihaooo/orange/models/couponModel"
	string2 "github.com/lilihaooo/orange/utils/str"
	"github.com/robfig/cron/v3"
)

// 优惠券发放

func couponIssuance() {
	/*
		逻辑:
		一. 每分钟查询一次符合条件的预发券信息 例如: x条符合的记录
			条件:
			1. 券状态为1
			2. 预发券状态为1
			3. 当前时间等于发券时间点
		二. 循环x go func(coupon_id int64, num int64) 发券(添加已发数量和写日志 用事务)
		问题:
			1. 在循环中操作数据库的压力是不是太大了
	*/

	c := cron.New()
	_, err := c.AddFunc("@every 1m", func() {
		err := couponModel.IssueCoupons(logWrite)
		if err != nil {
			logWrite.Error("发券脚本执行失败:" + err.Error())
		} else {
			logWrite.Info(string2.CurrentTimeYMDHIS() + "发券脚本执行成功")
		}
	})
	if err != nil {
		panic(err)
	}
	c.Start()
	select {}
}
