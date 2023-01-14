package baseHandler

import (
	"github.com/gin-gonic/gin"
	handlers "orange/handler"
	"orange/help"
	"orange/middleware/jwt"
	"orange/models"
	"orange/models/baseModel"
	"orange/utils/upload"
	"strconv"
)

type adminLogin struct {
	Username string `form:"username" json:"username" xml:"username"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// 用户登录
func Login(c *gin.Context) {
	var json adminLogin
	if err := c.ShouldBindJSON(&json); err != nil {
		handlers.FailWithParams(c, err)
		return
	}
	// 根据账号查找用户
	admin := baseModel.Admin{
		Username: json.Username,
	}

	admin, err := admin.GetAdmin()

	if err != nil {
		handlers.FailWithMessage(c, "用户不存在")
		return
	}

	// 根据密码盐验证密码是否正确
	if !baseModel.CheckPassword(json.Password, admin.Salt, admin.Password) {
		handlers.FailWithMessage(c, "密码错误")
		return
	}

	// 生成token
	token, err := jwt.GenerateToken(admin.ID, admin.Username, admin.Salt)
	ret := map[string]string{"token": token}
	handlers.Success(c, ret)
	return
}

// 用户退出
func Logout(c *gin.Context) {
	token := c.Request.Header.Get("X-Token")
	res := jwt.DestroyToken(token)
	if !res {
		handlers.FailWithParams(c, nil)
		return
	}
	handlers.Success(c, nil)
	return
}

func AdminList(c *gin.Context) {
	search := help.SearchParamsFormat(c)
	// 查询是否有username、name、mobile等查询参数
	search["username"] = c.Query("username")
	search["mobile"] = c.Query("mobile")
	search["is_deleted"] = c.Query("is_deleted")
	adminModel := &baseModel.Admin{}
	admin, total, err := adminModel.GetAdminList(search)
	if err != nil {
		handlers.FailWithSystem(c, err)
		return
	}
	response := make(map[string]interface{})
	response["list"] = admin
	response["total"] = total
	handlers.Success(c, response)

}

func AdminAdd(c *gin.Context) {
	admin := &baseModel.Admin{}
	err := c.BindJSON(admin)
	if err != nil {
		handlers.FailWithParams(c, err)
		return
	}

	err = admin.AddAdmin()

	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func AdminUpdate(c *gin.Context) {
	admin := &baseModel.Admin{}
	err := c.BindJSON(admin)
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	err = admin.UpdateAdmin()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, nil)
}

func AdminDelete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		handlers.FailWithParams(c, nil)
		return
	}
	idInt, _ := strconv.ParseInt(id, 10, 64)
	if idInt == 0 {
		handlers.FailWithParams(c, nil)
		return
	}

	admin := baseModel.Admin{
		Model: models.Model{ID: idInt},
	}
	admin, err := admin.GetAdmin()
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	//删除该用户redis中的数据
	jwt.DestroyTokenByUserInfo(admin.ID, admin.Username, admin.Salt)

	admin.DeleteAdmin()
	handlers.Success(c, nil)
}

// 将删除的管理员恢复
func AdminRecover(c *gin.Context) {
	admin := &baseModel.Admin{}
	err := c.BindJSON(admin)
	if err != nil {
		handlers.FailWithParams(c, err)
		return
	}
	admin.RecoverAdmin()
	handlers.Success(c, nil)
}

// 头像上传
func AdminUpload(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		handlers.FailWithParams(c, nil)
	}

	//上传图片
	path, err := upload.FileUpload(c, file, "admin")
	if err != nil {
		handlers.FailWithMessage(c, err.Error())
		return
	}
	handlers.Success(c, path)
}

// 查找管理员基本信息
func AdminInfo(c *gin.Context) {
	var userInfo map[string]interface{}
	userInfo = make(map[string]interface{})
	// 根据账号查找用户
	admin := baseModel.Admin{
		Username: c.GetString("username"),
	}
	admin, err := admin.GetAdmin()

	if err != nil {
		handlers.FailWithMessage(c, "用户不存在")
		return
	}
	userInfo["avatar"] = admin.Avatar
	handlers.Success(c, userInfo)
}
