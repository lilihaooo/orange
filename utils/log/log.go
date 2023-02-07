package log

import (
	"errors"
	"github.com/lilihaooo/orange/utils/file"
	string2 "github.com/lilihaooo/orange/utils/str"
	"os"
)

// LogFileWriter 根据日期获取新的日志文件
type LogFileWriter struct {
	File *os.File
	Date string
	Path string
}

// 返回日志io.Writer
func NewLogFileWriter(endpoint string, category string) *LogFileWriter {
	var log LogFileWriter
	var err error
	// 获取年月日
	log.Date = string2.CurrentTimeYMDHIS()
	// 根据endpoint来判断是前端api还是后端api
	switch endpoint {
	case "backend":
		log.Path = "storage/logs/admin/" + log.Date + "/"
	case "frontend":
		log.Path = "storage/logs/api/" + log.Date + "/"
	case "crontab":
		log.Path = "storage/logs/crontab/" + log.Date + "/"
	case "wechatPay":
		log.Path = "storage/logs/wechatPay/" + log.Date + "/"
	default:
		log.Path = "storage/logs/system/" + log.Date + "/"
	}
	fileName := category + ".api-log"
	// 判断该日期是否有文件存在 没有就创建
	if !file.CheckExist(log.Path + fileName) {
		_, err = file.MustOpen(fileName, log.Path)
		if err != nil {
			panic("创建日志文件失败:" + err.Error())
		}
	}
	if log.File, err = os.OpenFile(log.Path+fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600); err != nil {
		panic("创建日志文件失败:" + err.Error())
	}
	return &log
}

func (p *LogFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.File == nil {
		return 0, errors.New("file not opened")
	}
	// 判断该日期是否有文件存在 没有就创建
	if !file.CheckExist(p.File.Name()) {
		// 将p.File.Name()
		_, err = file.MustOpen(p.File.Name()[len(p.Path):], p.Path)
		if err != nil {
			panic("创建日志文件失败:" + err.Error())
		}
		p.File, _ = os.OpenFile(p.File.Name(), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	}
	n, e := p.File.Write(data)

	return n, e
}
