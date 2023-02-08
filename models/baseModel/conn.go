package baseModel

import (
	"github.com/lilihaooo/orange/db/conn/mysql"
	"gorm.io/gorm"
)

var conn *gorm.DB

func InitConn() {
	// 构建base数据库
	conn = mysql.GetConn()
}
