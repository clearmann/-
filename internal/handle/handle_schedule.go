package handle

import (
	g "backend/internal/global"
	"backend/internal/model"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type Schedule struct {
}
type AddScheduleRequest struct {
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content" binding:"required"`
	ScheduleTime string `json:"scheduleTime" binding:"required"`
	Count        int    `json:"count"`
	RepeatTimes  int    `json:"repeatTimes"`
	Style        int    `json:"style"`
}

var (
	IsInitEmailScheduler bool = false
)

func (*Schedule) AddSchedule(c *gin.Context) {
	var addScheduleRequest AddScheduleRequest
	if err := c.ShouldBindJSON(&addScheduleRequest); err != nil {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "格式错误",
		})
		return
	}
	if addScheduleRequest.Title == "" || addScheduleRequest.Content == "" || addScheduleRequest.ScheduleTime == "" {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "格式错误",
		})
		return
	}
	db := GetDB(c)
	session := sessions.Default(c)
	userID := session.Get(g.CTX_USER_AUTH).(int)
	log.Println("userID", userID)
	email, _ := model.GetEmailByUserID(db, userID)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	parsedTime, _ := time.ParseInLocation("2006-01-02 15:04:05", addScheduleRequest.ScheduleTime, loc)
	scheduleID, err := model.CreateSchedule(db, addScheduleRequest.Title, addScheduleRequest.Content, parsedTime, addScheduleRequest.RepeatTimes, addScheduleRequest.Style, userID, email)
	if err != nil {
		log.Println("创建日程失败...")
		return
	}
	go scheduleReminderNotification(db, scheduleID)
	if !IsInitEmailScheduler {
		go InitEmailScheduler(c)
		IsInitEmailScheduler = true
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
}
func scheduleReminderNotification(db *gorm.DB, scheduleID int) {
	schedule, err := model.GetScheduleByScheduleID(db, scheduleID)
	if err != nil {
		log.Println("获取日程失败...")
		return
	}
	ScheduleMap.Store(schedule.ID, schedule)
}
func (*Schedule) DelSchedule(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.PostForm("scheduleID"))
	db := GetDB(c)
	session := sessions.Default(c)
	userID := session.Get(g.CTX_USER_AUTH).(int)
	log.Println(userID, scheduleID)
	if err := model.DelSchedule(db, userID, scheduleID); err != nil {
		log.Println(err.Error())
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
	ScheduleMap.Delete(scheduleID)
}

type UpdateScheduleRequest struct {
	ScheduleID   int    `json:"scheduleID"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ScheduleTime string `json:"scheduleTime"`
	Style        int    `json:"style"`       // 1 为一次性日程  2 为重复日程  3 为周一至周五 日程
	RepeatTimes  int    `json:"repeatTimes"` //每五分钟发送一次 要发送几次
	Count        int    `json:"count"`
}

func (*Schedule) UpdateSchedule(c *gin.Context) {
	var updateScheduleRequest UpdateScheduleRequest
	if err := c.ShouldBindJSON(&updateScheduleRequest); err != nil {
		log.Println(err.Error())
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "请求参数格式错误",
		})
		return
	}
	db := GetDB(c)
	session := sessions.Default(c)
	userID := session.Get(g.CTX_USER_AUTH).(int)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	parsedTime, _ := time.ParseInLocation("2006-01-02 15:04:05", updateScheduleRequest.ScheduleTime, loc)
	var schedule = &model.Schedule{
		BaseModel:    model.BaseModel{ID: updateScheduleRequest.ScheduleID},
		UserID:       userID,
		Title:        updateScheduleRequest.Title,
		Content:      updateScheduleRequest.Content,
		ScheduleTime: parsedTime,
		Style:        updateScheduleRequest.Style,
		RepeatTimes:  updateScheduleRequest.RepeatTimes,
	}
	log.Println(schedule.UserID, schedule.ID)
	if err := model.UpdateSchedule(db, schedule); err != nil {
		c.JSON(200, gin.H{
			"code":    200,
			"message": "error",
			"data":    "更新失败...",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
	})
	ScheduleMap.Store(updateScheduleRequest.ScheduleID, schedule)
}

type ScheduleResponse struct {
	ID           int    `json:"ID"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	ScheduleTime string `json:"scheduleTime"`
	Style        int    `json:"style"`                       // 1 为一次性日程  2 为重复日程  3 为周一至周五 日程
	RepeatTimes  int    `json:"repeatTimes" gorm:"type:int"` //每五分钟发送一次 要发送几次
}

func (*Schedule) GetScheduleList(c *gin.Context) {
	db := GetDB(c)
	session := sessions.Default(c)
	UserID := session.Get(g.CTX_USER_AUTH).(int)

	total, schedules, err := model.GetScheduleListByUserID(db, UserID)
	if err != nil {
		log.Println("查询Schedule失败...")
	}
	data := make([]ScheduleResponse, total)
	for i, schedule := range *schedules {
		data[i].ID = schedule.ID
		data[i].Title = schedule.Title
		data[i].Content = schedule.Content
		data[i].ScheduleTime = schedule.ScheduleTime.Format("2006-01-02 15:04")
		data[i].Style = schedule.Style
		data[i].RepeatTimes = schedule.RepeatTimes
	}
	fmt.Println(data)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}
