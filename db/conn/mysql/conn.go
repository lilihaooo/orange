package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"orange/settings"
	"time"
)

type Conn = *gorm.DB
type DefaultConn = Conn

// 数据控制结构体
type DBManager struct {
	DefaultConn
	Connections map[string]Conn
}

var DB *DBManager

func Connect(settings *settings.MySQLConfig) *DBManager {
	db := &DBManager{}
	conn := newConn(settings)
	// todo 最大连接要设置成mysql最大时长的一半
	conn.DB().SetConnMaxLifetime(time.Minute)
	db.DefaultConn = conn
	db.DefaultConn.LogMode(true)
	//禁用复数形式
	db.SingularTable(true)
	DB = db
	return db
}

func newConn(config *settings.MySQLConfig) Conn {
	args := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&%s",
		config.User, config.Password, config.Host, config.Port, config.DB, config.Options)
	db, err := gorm.Open("mysql", args)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetConn() Conn {
	return DB.DefaultConn
}
