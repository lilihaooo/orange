package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	myRedis "github.com/lilihaooo/orange/db/conn/redis"
	string2 "github.com/lilihaooo/orange/utils/str"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 加密
var jwtSecret = []byte("lihao")

// 发行人
var issuer = "lihao"

// 过期时间
var expire time.Duration = 24 * 30

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(id int64, username string, Salt string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(expire * time.Hour)

	claims := Claims{
		id,
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    issuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	//使用redis储存token,
	redisClient := myRedis.Pool.Get()
	defer redisClient.Close()
	viewToken := makeViewToken(id, username, Salt)
	redisClient.Do("SET", myRedis.KeyPrefix+viewToken, token)
	// 设置用户有效期
	redisClient.Do("EXPIRE", myRedis.KeyPrefix+viewToken, expire*3600)
	return viewToken, err
}

// 生成用户token
func makeViewToken(id int64, username string, salt string) string {
	//key就是根据用户信息组合的特殊字符串 后md5加密(作为新的token)  value就是token
	//todo key使用id和用户名加密
	//key := "authorization_" + strconv.FormatInt(id, 10) + "_" + "_" + help.CurrentTimeYMDHIS()
	key := "authorization_" + strconv.FormatInt(id, 10) + "_" + username + "_" + salt
	//key : 前缀加加密后的key
	return string2.EncodeMD5(key)
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// DestroyToken 删除根据token删除redis中数据
func DestroyToken(token string) bool {
	key := myRedis.KeyPrefix + token
	redisClient := myRedis.Pool.Get()
	defer redisClient.Close()
	_, err := redisClient.Do("DEL", key)
	if err != nil {
		return false
	}
	return true
}

// 根据用户信息删除redis中的数据
func DestroyTokenByUserInfo(id int64, username string, salt string) bool {
	viewToken := makeViewToken(id, username, salt)
	key := myRedis.KeyPrefix + viewToken
	redisClient := myRedis.Pool.Get()
	defer redisClient.Close()
	_, err := redisClient.Do("DEL", key)
	if err != nil {
		return false
	}
	return true
}

func JWT() gin.HandlerFunc {

	return func(c *gin.Context) {
		var code int

		code = 20000
		// 查询参数token
		token := c.DefaultQuery("token", "")
		// 如果不存在 从header头读取信息
		if token == "" {
			// 读取Authorization
			authorization := c.Request.Header.Get("Authorization")
			if authorization != "" {
				bearer := strings.Split(authorization, " ")
				if len(bearer) < 2 {
					code = 50008
					c.JSON(http.StatusOK, gin.H{
						"code":    code,
						"message": "用户验证失败",
					})
					c.Abort()
					return
				}

				token = bearer[1]
			}
			// 读取Authorization
			xToken := c.Request.Header.Get("X-Token")
			if xToken != "" {
				token = xToken
			}
		}

		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "用户验证失败",
			})
			c.Abort()
			return
		}

		//-------------------------------------
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

		//-------------------------------------
		// 解析token
		claims, err := ParseToken(token)
		if err != nil {
			code = 50008
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = 50014
		}
		if code != 20000 {
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": "用户验证失败",
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("token", token)
		c.Set("userId", claims.ID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
