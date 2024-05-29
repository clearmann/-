package middleware

import (
	g "backend/internal/global"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

// WithGormDB 将 gorm.DB 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_DB).(*gorm.DB) 来使用
func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(g.CTX_DB, db)
		ctx.Next()
	}
}

// CORS 跨域请求
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	})
}

// WithRedisDB 将 redis.Client 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_RDB).(*redis.Client) 来使用
func WithRedisDB(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(g.CTX_RDB, rdb)
		ctx.Next()
	}
}

// WithCookieStore 基于 cookie 的 session
func WithCookieStore(name, secret string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{Path: "/", MaxAge: g.Conf.Session.MaxAge})
	return sessions.Sessions(name, store)
}
