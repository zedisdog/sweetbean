package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/zedisdog/sweetbean/tools"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

//GenAuthMiddleware generate an auth middleware
//	claims: the customer Claims object
//	key: the key for validate sign
//	isUserExists: a function to determine if account is exists
//Deprecated: use BuildAuth instead
func GenAuthMiddleware(key string, isUserExists func(id interface{}) bool) func(*gin.Context) {
	return func(c *gin.Context) {
		var token string
		if c.Request.Header.Get("Authorization") != "" {
			arr := strings.Split(c.Request.Header.Get("Authorization"), " ")
			if len(arr) < 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问1"})
				return
			}
			token = arr[1]
		} else if c.Query("token") != "" {
			token = c.Query("token")
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问2"})
			return
		}

		t, err := tools.Parse(token, []byte(key))
		if err != nil || !t.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问3"})
			return
		}

		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			var id interface{}
			id, ok = claims["jti"]
			if !ok || !isUserExists(id) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问4"})
				return
			}
			c.Set("id", id)
			c.Set("claims", claims)
		}

		c.Next()
	}
}

//BuildAuth
//  通过加密用的key和查找用户是否存在的函数构造身份验证中间件，中间件通过header中的Authorization字段或者url中queryString的token字段来获取token
func BuildAuth(key string, isUserExists func(id uint64) bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var token string
		if ctx.Request.Header.Get("Authorization") != "" {
			arr := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
			if len(arr) < 2 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问1"})
				return
			}
			token = arr[1]
		} else if ctx.Query("token") != "" {
			token = ctx.Query("token")
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问2"})
			return
		}

		t, err := tools.Parse(token, []byte(key))
		if err != nil || !t.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问3"})
			return
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "未授权的访问7",
			})
			return
		}

		var jti interface{}
		jti, ok = claims["jti"]
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "未授权的访问4",
			})
			return
		}

		var id uint64
		id, err = strconv.ParseUint(jti.(string), 10, 64)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "未授权的访问5",
			})
			return
		}

		if !isUserExists(id) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "未授权的访问6",
			})
			return
		}

		ctx.Set("id", id)
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
