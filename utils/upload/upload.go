package upload

import (
	"github.com/gin-gonic/gin"
	string2 "github.com/lilihaooo/orange/utils/str"
	"mime/multipart"
	"os"
)

func FileUpload(c *gin.Context, file *multipart.FileHeader, who string) (string, error) {
	date := string2.CurrentTimeYMDHIS()
	dir := "./resource/" + who + "/" + date + "/"
	//检测目录是否存在, 不存在就创建
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	path := dir + file.Filename
	return path, c.SaveUploadedFile(file, path)
}
