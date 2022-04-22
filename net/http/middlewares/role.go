package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CheckRole interface {
	IsRole(id interface{}, roleName string) bool
}

func GenRoleMiddleware(svc CheckRole, roleName string) func(*gin.Context) {
	return func(c *gin.Context) {
		if !svc.IsRole(c.MustGet("id"), roleName) {
			c.AbortWithStatusJSON(http.StatusForbidden, map[string]string{"message": "未授权的访问"})
			return
		}
		c.Next()
	}
}
