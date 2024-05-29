package middleware

import (
	g "backend/internal/global"
	jwt "backend/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

// JWTAuth 基于 JWT 的授权
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("中间件验证jwt...")
		// 从 URL 中获取 token
		token := c.Request.URL.Query().Get("token")
		// 从请求头中获取 token
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			log.Println("请求头中不含有 Authorization")
		}
		// token 的正确格式为 "Bearer [tokenString]"
		parts := strings.Split(authorization, " ")
		if token == "" {
			token = parts[1]
		}
		if token == "" {
			c.JSON(200, gin.H{
				"code":    200,
				"message": "error",
				"data":    "身份验证失败",
			})
			return
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(200, gin.H{
				"code":    200,
				"message": "error",
				"data":    "TOKEN 不正确，请重新登录",
			})
			return
		}
		// 判断token 是否过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			c.JSON(200, gin.H{
				"code":    200,
				"message": "error",
				"data":    "TOKEN 已过期，请重新登陆",
			})
			return
		}
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, claims.UserID)
		err = session.Save()
		if err != nil {
			return
		}
		log.Println("身份验证成功...")
		c.Next()
	}
}

//func PermissionCheck() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		if c.GetBool("skip_check") {
//			c.Next()
//			return
//		}
//		//db := c.MustGet(g.CTX_DB).(*gorm.DB)
//		auth, err := handle.GetCurrentUserAuth(c)
//		if err != nil {
//			handle.ReturnError(c, g.ErrDbOp, nil)
//		}
//		if auth.IsSuper {
//			zap.L().Info("[middleware-PermissionCheck]: super admin no need to check, pass!")
//			c.Next()
//			return
//		}
//		url := c.FullPath()[4:]
//		method := c.Request.Method
//		zap.L().Info(fmt.Sprintf("[middleware-PermissionCheck] %v, %v, %v\n", auth.Username, url, method))
//		//for _, role := range auth.Roles {
//		//}
//		// 目前先全部通过
//		c.Next()
//	}
//}
