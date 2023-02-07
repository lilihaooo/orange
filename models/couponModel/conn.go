package couponModel

import (
	"github.com/lilihaooo/orange/db/conn/mysql"
)

var conn mysql.Conn

func InitConn() {
	// 构建coupon数据库
	conn = mysql.GetConn()
}
