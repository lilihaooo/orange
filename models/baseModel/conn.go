package baseModel

import (
	"github.com/lilihaooo/orange/db/conn/mysql"
)

var conn mysql.Conn

func InitConn() {
	// 构建base数据库
	conn = mysql.GetConn()
}
