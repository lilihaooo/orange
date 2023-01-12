package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	myRedis "orange/db/conn/redis"
	"orange/help"
	"strconv"
	"time"
)

// var apiJwtSecret = []byte("zhouqi")
// var apiIssuer = "zhouqi"
var apiExpire time.Duration = 24 * 30

type ApiClaims struct {
	ID     int64  `json:"id"`
	Mobile string `json:"mobile"`
	jwt.StandardClaims
}

func ApiGenerateToken(id int64, mobile string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(apiExpire * time.Hour)

	claims := ApiClaims{
		id,
		mobile,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	// 使用redis存储token 参数token为md5
	redisClient := myRedis.Pool.Get()
	defer redisClient.Close()
	key := "authorization_" + strconv.FormatInt(id, 10) + "_" + mobile + "_" + help.CurrentTimeYMDHIS()
	paramsToken := help.EncodeMD5(key)
	redisClient.Do("SET", myRedis.KeyPrefix+paramsToken, token)
	// 设置用户有效期
	redisClient.Do("EXPIRE", myRedis.KeyPrefix+paramsToken, apiExpire*3600)
	return paramsToken, err
}

func ApiParseToken(token string) (*ApiClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &ApiClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*ApiClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func ApiDestroyToken(token string) bool {

	return true
}

func ApiJWT() gin.HandlerFunc {

	return func(c *gin.Context) {
		var code int

		code = 20000
		token := ""
		// 查询参数token 兼容url GET传参和header传参
		token = c.Request.Header.Get("Token")
		token2 := c.DefaultQuery("token", "")
		if token == "" && token2 == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    50001,
				"message": "登陆异常,请重新登陆",
			})
			c.Abort()
			return
		}
		if token2 != "" && token == "" {
			token = token2
		}
		// 使用redis读取token 得到真正的jwt
		redisClient := myRedis.Pool.Get()
		defer redisClient.Close()

		token, err := redis.String(redisClient.Do("GET", myRedis.KeyPrefix+token))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    50014,
				"message": "登陆异常,请重新登陆",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := ApiParseToken(token)
		if err != nil {
			code = 50008
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = 50014
		}
		if code != 20000 {
			c.JSON(http.StatusOK, gin.H{
				"code":    50014,
				"message": "用户验证失败",
			})
			c.Abort()
			return
		}

		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("token", token)
		c.Set("mobile", claims.Mobile)
		c.Set("user_id", claims.ID)
		c.Next()
	}
}
