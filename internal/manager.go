package internal

import (
	"backend/internal/handle"
	"backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

var (
	userAuthAPI handle.UserAuth // 用户验证
	userInfoAPI handle.UserInfo // 用户信息
	scheduleAPI handle.Schedule
)

func RegisterHandlers(r *gin.Engine) {
	registerBaseHandlers(r)
}

// 通用接口: 全部不需要 登录 + 鉴权
func registerBaseHandlers(r *gin.Engine) {
	base := r.Group("/api/web")
	base.POST("/user/account/register/", userAuthAPI.Register)
	base.POST("/user/account/token/", userAuthAPI.Token)
	base.POST("/user/account/info/", userInfoAPI.Info)
	protected := base.Group("/user")
	{
		protected.Use(middleware.JWTAuth())
		protected.POST("/schedule/", scheduleAPI.AddSchedule)    //添加一个日程
		protected.GET("/schedule/", scheduleAPI.GetScheduleList) //查询所有日程
		protected.DELETE("/schedule/", scheduleAPI.DelSchedule)  //删除一个日程
		protected.PUT("/schedule/", scheduleAPI.UpdateSchedule)  //修改一个日程
	}
}
