package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckRole interface {
	IsRole(id interface{}, roleName string) bool
}

// GenRoleMiddleware
// Deprecated: use GenChecker instead
func GenRoleMiddleware(svc CheckRole, roleNames ...string) func(*gin.Context) {
	return func(c *gin.Context) {
		for _, roleName := range roleNames {
			if svc.IsRole(c.MustGet("id"), roleName) {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, map[string]string{"message": "未授权的访问"})
	}
}
