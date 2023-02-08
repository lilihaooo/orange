package validCheck

import (
	"strconv"
	"strings"
	"time"
)

// StrISMs 判断字符串是否是时分格式的
func StrISMs(str string) bool {
	if len(str) != 5 {
		return false
	}
	Ms := strings.Split(str, ":")
	if len(Ms) != 2 {
		return false
	}
	m, err := strconv.Atoi(Ms[0])
	if err != nil {
		return false
	}
	s, err := strconv.Atoi(Ms[1])
	if err != nil {
		return false
	}
	if m < 0 || m >= 24 {
		return false
	}
	if s < 0 || s >= 60 {
		return false
	}
	return true
}

// StrISMs 判断字符串是否是时分格式的
func StrISYmd(str string) bool {
	_, err := time.Parse("2006-01-02", str)
	if err != nil {
		return false
	}
	return true
}
