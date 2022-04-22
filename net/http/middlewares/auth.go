package middlewares

import (
	"github.com/zedisdog/sweetbean/tools"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

//GenAuthMiddleware generate an auth middleware
//	claims: the customer Claims object
//	key: the key for validate sign
//	isUserExists: a function to determine if account is exists
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
