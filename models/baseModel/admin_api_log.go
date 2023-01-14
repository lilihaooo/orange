package baseModel

import (
	"time"
)

type AdminApiLog struct {
	ID        int64      `gorm:"primary_key"json:"id"`
	UserId    int64      `json:"user_id"`
	Username  string     `json:"username"`
	Role      string     `json:"role"`
	Host      string     `json:"host"`
	Path      string     `json:"path"`
	Method    string     `json:"method"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at"`
}

func (AdminApiLog) TableName() string {
	return "admin_api_log"
}

// todo  为什么这里要用一个函数返回 不直接用
func NewAdminApiLogModel() AdminApiLog {
	return AdminApiLog{}
}

// 查询公告
func (m *AdminApiLog) GetLogList(params map[string]interface{}) (logs []AdminApiLog, count int, err error) {
	db := conn.Model(&AdminApiLog{})
	// start_time
	if start_time, ok := params["start_time"]; ok && start_time.(string) != "" {
		db = db.Where("created_at >= ? ", start_time.(string)+" 00:00:00")
	}

	// end_time
	if end_time, ok := params["end_time"]; ok && end_time.(string) != "" {
		db = db.Where("created_at <= ? ", end_time.(string)+" 23:59:59")
	}
	// todo count是否应该放在where 后面
	db.Count(&count)

	err = db.Offset((params["page"].(int) - 1) * params["pageSize"].(int)).Limit(params["pageSize"]).Order("id DESC").Find(&logs).Error
	if err != nil {
		return
	}
	return
}

// 添加日志
func (m *AdminApiLog) AddAdminApiLog() {
	conn.Create(m)
}

type result struct {
	Username string `json:"username"`
	Date     string `json:"date"`
	Count    int    `json:"count"`
}

// 统计
func (m *AdminApiLog) StatLog() (*[]result, error) {
	var res []result
	db := conn.Model(&AdminApiLog{})
	err := db.Select("username, DATE_FORMAT(created_at, '%d/%m/%y') as date, count(1) as count").Where("method = ?", "delete").Group("username, date").Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}
