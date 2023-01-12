package baseModel

import (
	"errors"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-ozzo/ozzo-validation/v4"
	"orange/help"
	"orange/models"
	"regexp"
	"strconv"
	"strings"
)

type Admin struct {
	models.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

// 去除空格符号
func (m *Admin) trimField() error {
	m.Username = strings.Trim(m.Username, " ")
	m.Password = strings.Trim(m.Password, " ")
	m.Mobile = strings.Trim(m.Mobile, " ")
	if len(m.Mobile) != 11 {
		return errors.New("请输入11位手机号")
	}

	// 正则表达式匹配字母或数字 不允许出现汉字
	pattern := `^[A-Za-z0-9]{6,24}$`
	matched, _ := regexp.MatchString(pattern, m.Username)
	if !matched {
		return errors.New("用户名必须为6-24位字母或数字")
	}

	if (m.ID != 0 && m.Password != "") || m.ID == 0 {
		pattern = `^[_a-z0-9-]{6}`
		matched, _ = regexp.MatchString(pattern, m.Password)
		if !matched {
			return errors.New("密码必须为至少6位字母或数字")
		}
	}

	return nil
}

// 验证
func (m Admin) Validate() error {
	return validation.ValidateStruct(&m,
		// 名称不得为空,且大小为1-50字
		validation.Field(
			&m.Username,
			validation.Required.Error("用户名不得为空"),
			validation.Length(1, 50).Error("名称为1-50个字母或数字")),
		// 密码不得为空,且大于6字
		validation.Field(
			&m.Password,
			validation.Required.Error("密码不得为空"),
			validation.Length(6, 0).Error("密码不得小于6位")),
		// 电话号码必须为11位
		validation.Field(
			&m.Mobile,
			validation.Required.Error("电话号码不得为空"),
			validation.Length(11, 11).Error("请输入11位手机号码")),
		// 邮箱格式
		validation.Field(
			&m.Email,
			is.Email.Error("邮箱格式错误")),
		// 头像地址
		validation.Field(
			&m.Avatar,
			is.URL.Error("头像地址格式错误")),
	)
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
	err = db.Offset((params["page"].(int) - 1) * params["pageSize"].(int)).Limit(params["pageSize"]).Find(&admin).Error
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
	if err := m.trimField(); err != nil {
		return err
	}
	err = m.Validate()
	if err != nil {
		return err
	}
	// 用户必须唯一
	if !conn.Where("username = ?", m.Username).First(&m).RecordNotFound() {
		return errors.New("用户名已存在")
	}

	// 生成密码盐
	m.Salt = help.EncodeMD5(help.RandString(10))
	// 对密码进行加密
	m.Password = help.EncodeMD5(m.Password + m.Salt)
	if err := conn.Create(m).Error; err != nil {
		return err
	}
	return nil

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
	if err := m.trimField(); err != nil {
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
		adminDb.Password = help.EncodeMD5(m.Password + adminDb.Salt)
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
	if dbpassword == help.EncodeMD5(password+salt) {
		return true
	}
	return false
}
