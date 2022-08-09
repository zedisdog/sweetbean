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

//BuildAuth 通过加密用的key和查找用户是否存在的函数构造身份验证中间件，中间件通过header中的Authorization字段或者url中queryString的token字段来获取token
//Deprecated: use middlewares.auth generate by middlewares.NewAuth instead.
func BuildAuth(key string, isUserExists func(id uint64) (bool, error)) func(ctx *gin.Context) {
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

		exists, err := isUserExists(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		if !exists {
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

func NewAuthBuilder() *authBuilder {
	return &authBuilder{}
}

// authBuilder auth middleware builder. it parse token with given conditions.
type authBuilder struct {
	userIdentityFrom string                        //field name of user identity in token
	tokenIDFrom      string                        //filed name of token identity in token
	roleFrom         string                        //filed name of role name in token
	userExists       func(id uint64) (bool, error) //function to determine if user exists
	authKey          []byte                        //salt used by generate jwt signature
	cacheClaims      bool                          //if cache claims into context
}

func (ab *authBuilder) WithUserIdentityFrom(jwtField string) *authBuilder {
	ab.userIdentityFrom = jwtField
	return ab
}

func (ab *authBuilder) WithTokenIDFrom(jwtField string) *authBuilder {
	ab.tokenIDFrom = jwtField
	return ab
}

func (ab *authBuilder) WithRoleFrom(jwtField string) *authBuilder {
	ab.roleFrom = jwtField
	return ab
}

func (ab *authBuilder) WithUserExistsFunc(f func(id uint64) (bool, error)) *authBuilder {
	ab.userExists = f
	return ab
}

func (ab *authBuilder) WithAuthKey(key string) *authBuilder {
	ab.authKey = []byte(key)
	return ab
}

func (ab *authBuilder) WithClaimsCache() *authBuilder {
	ab.cacheClaims = true
	return ab
}

func (ab *authBuilder) Build() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var token string
		if ctx.Request.Header.Get("Authorization") != "" {
			arr := strings.Split(ctx.Request.Header.Get("Authorization"), " ")
			if len(arr) < 2 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "bearer token is invalid"})
				return
			}
			token = arr[1]
		} else if ctx.Query("token") != "" {
			token = ctx.Query("token")
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "no token found"})
			return
		}

		t, err := tools.Parse(token, ab.authKey)
		if err != nil || !t.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "token is invalid1"})
			return
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token is invalid2",
			})
			return
		}

		if ab.userIdentityFrom != "" {
			var IDStr interface{}
			IDStr, ok = claims[ab.userIdentityFrom]
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is invalid3",
				})
				return
			}

			var id uint64
			id, err = strconv.ParseUint(IDStr.(string), 10, 64)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is invalid4",
				})
				return
			}

			if ab.userExists != nil {
				exists, err := ab.userExists(id)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"message": "token is invalid5",
					})
				}
				if !exists {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": "token is invalid6",
					})
					return
				}
			}

			ctx.Set("id", id)
		}

		if ab.roleFrom != "" {
			var role interface{}
			role, ok = claims[ab.roleFrom]
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token is invalid7",
				})
				return
			}
			ctx.Set("role", role.(string))
		}

		if ab.cacheClaims {
			ctx.Set("claims", claims)
		}

		ctx.Next()
	}
}
