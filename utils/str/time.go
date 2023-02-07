package string

import (
	"time"
)

// CurrentTimeYMD 获取年月日
func CurrentTimeYMD() string {
	return time.Now().Format("2006-01-02")
}

// CurrentTimeYMDHIS 获取年月日时分秒
func CurrentTimeYMDHIS() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// CurrentTimeHI 获取时分
func CurrentTimeHI() string {
	return time.Now().Format("15:04")
}
