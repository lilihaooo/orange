package baseModel

import (
	"errors"
	"github.com/lilihaooo/orange/models"
	string2 "github.com/lilihaooo/orange/utils/str"
	"github.com/lilihaooo/orange/utils/validCheck"
	"strconv"
)

// 使用 validator验证包
type Admin struct {
	models.Model
	Username string `json:"username" validate:"required,alphanum,min=4,max=8" label:"用户名"`
	Password string `json:"password" validate:"required,min=6,max=8"`
	Salt     string `json:"salt"`
	Mobile   string `json:"mobile" validate:"len=11" label:"电话"`
	Avatar   string `json:"avatar"` //todo 头像格式未验证
	Email    string `json:"email" validate:"email"`
}

// 查询管理员
func (m *Admin) GetAdminList(params map[string]interface{}) (admin []Admin, count int, err error) {
	//Unscoped()  查询结果包含被软删除的
	db := conn.Model(&Admin{}).Unscoped()
	// username
	if username, ok := params["username"]; ok && username.(string) != "" {
		db = db.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	// mobile
	if mobile, ok := params["mobile"]; ok && mobile.(string) != "" {
		db = db.Where("mobile LIKE ?", "%"+mobile.(string)+"%")
	}
	// 是否删除
	if isDeleted, ok := params["is_deleted"]; ok {
		if isDeleted.(string) != "" {
			is, _ := strconv.Atoi(isDeleted.(string))
			// 查询已经删除的用户
			if is == 1 {
				db = db.Where("deleted_at IS NOT NULL")
			} else { // 未删除的用户
				db = db.Where("deleted_at IS NULL")
			}
		}
	}
	db.Count(&count)
	err = db.Offset((params["search"].(int) - 1) * params["pageSize"].(int)).Limit(params["pageSize"]).Find(&admin).Error
	if err != nil {
		return
	}
	return
}

// 查询所有管理员
func GetAdminAll() (admin []Admin, err error) {
	err = conn.Model(&Admin{}).Find(&admin).Error
	if err != nil {
		return
	}
	return
}

// 查询管理员
func (m *Admin) GetAdmin() (admin Admin, err error) {
	err = conn.Where(m).First(&admin).Error
	if err != nil {
		return Admin{}, err
	}
	return admin, nil
}

// 添加管理员
func (m *Admin) AddAdmin() (err error) {
	err = validCheck.Validate(m)
	if err != nil {
		return err
	}
	// 用户必须唯一
	if !conn.Where("username = ?", m.Username).First(&m).RecordNotFound() {
		return errors.New("用户名已存在")
	}
	// 手机号是否可以重复
	// 生成密码盐
	m.Salt = string2.EncodeMD5(string2.RandString(10))
	// 对密码进行加密
	m.Password = string2.EncodeMD5(m.Password + m.Salt)
	return conn.Create(m).Error
}

// 删除管理员
func (m *Admin) DeleteAdmin() {
	conn.Delete(m)
}

// 恢复管理员
func (m *Admin) RecoverAdmin() {
	conn.Model(Admin{}).Unscoped().Where(" id  = ?", m.ID).UpdateColumn("deleted_at", nil)
}

// 更新管理员
func (m *Admin) UpdateAdmin() (err error) {
	err = validCheck.Validate(m)
	if err != nil {
		return err
	}
	if m.ID == 0 {
		return errors.New("ID不等为空")
	}
	var adminDb Admin
	res := conn.Unscoped().Where("id = ?", m.ID).First(&adminDb)
	if res.RowsAffected == 0 {
		return errors.New("管理员不存在")
	}
	if m.Password != "" {
		// 对密码进行加密
		adminDb.Password = string2.EncodeMD5(m.Password + adminDb.Salt)
	}

	//之所以要赋值给新的结构体变量是因为只修改一部分字段,而不是全部字段1
	adminDb.Mobile = m.Mobile
	adminDb.Email = m.Email
	adminDb.Avatar = m.Avatar
	//todo update
	if err := conn.Unscoped().Save(&adminDb).Error; err != nil {
		return err
	}
	return nil

}

// 检测密码
func CheckPassword(password, salt, dbpassword string) bool {
	if dbpassword == string2.EncodeMD5(password+salt) {
		return true
	}
	return false
}
