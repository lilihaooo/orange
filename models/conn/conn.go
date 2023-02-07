package conn

import (
	"github.com/lilihaooo/orange/models/baseModel"
	"github.com/lilihaooo/orange/models/couponModel"
)

func InitConn() {
	// todo 这一层的目的是什么, 为什么不全局使用一个数据库实例?
	// 启动基础数据库
	baseModel.InitConn()

	// 启动用户数据库
	couponModel.InitConn()

	// 产品数据库
	//productModel.InitConn()

}
