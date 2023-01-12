package baseModel

import (
	"orange/db/conn/mysql"
)

var conn mysql.Conn

func InitConn() {
	// 构建base数据库
	conn = mysql.GetConn()
	conn.AutoMigrate(&Admin{}).AutoMigrate(&AdminApiLog{})
}
