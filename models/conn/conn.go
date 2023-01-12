package conn

import "orange/models/baseModel"

func InitConn() {
	// 启动基础数据库
	baseModel.InitConn()

	// 启动用户数据库
	//userModel.InitConn()

	// 产品数据库
	//productModel.InitConn()

}
