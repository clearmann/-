package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	BaseModel
	UserID       int       `json:"userID"`
	User         User      `gorm:"foreignKey:UserID"`
	Title        string    `json:"title" gorm:"type:varchar(256)"`
	Content      string    `json:"content" gorm:"type:varchar(256)"`
	ScheduleTime time.Time `json:"scheduleTime"`
	Style        int       `json:"style"`                       // 1 为一次性日程  2 为重复日程  3 为周一至周五 日程
	RepeatTimes  int       `json:"repeatTimes" gorm:"type:int"` //每五分钟发送一次 要发送几次
	Email        string    `json:"email"`
}

func CreateSchedule(db *gorm.DB, title string, content string, scheduleTime time.Time, repeatTimes, style, userID int, email string) (int, error) {
	var schedule = &Schedule{
		UserID:       userID,
		Title:        title,
		Content:      content,
		ScheduleTime: scheduleTime,
		RepeatTimes:  repeatTimes,
		Style:        style,
		Email:        email,
	}

	result := db.Create(&schedule)
	return schedule.ID, result.Error
}
func DelSchedule(db *gorm.DB, UserID, ScheduleID int) error {
	var schedule Schedule
	if err := db.Where("id = ?", ScheduleID).First(&schedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("未找到该日程")
		}
	}
	if schedule.UserID != UserID {
		return errors.New("权限不够，无法删除该日程")
	}
	if err := db.Delete(&schedule).Error; err != nil {
		return err
	}
	return nil
}
func GetScheduleListByUserID(db *gorm.DB, UserID int) (int64, *[]Schedule, error) {
	var bots *[]Schedule
	db = db.Model(&Schedule{})
	db = db.Where("user_id = ?", UserID)
	result := db.Find(&bots)
	total := result.RowsAffected
	return total, bots, result.Error
}
func GetScheduleByScheduleID(db *gorm.DB, scheduleID int) (*Schedule, error) {
	var schedule *Schedule
	db = db.Model(&Schedule{})
	db = db.Where("id = ?", scheduleID)
	result := db.Find(&schedule)
	return schedule, result.Error
}
func GetEmailByScheduleID(db *gorm.DB, scheduleID int) (string, error) {
	var schedule *Schedule
	db = db.Model(&Schedule{})
	db = db.Where("id = ?", scheduleID)
	result := db.Find(&schedule)
	return schedule.Title, result.Error
}
func UpdateSchedule(db *gorm.DB, request *Schedule) error {
	db = db.Model(&Schedule{}).Where("user_id = ? and id = ?", request.UserID, request.ID).
		Updates(Schedule{
			Title:        request.Title,
			Content:      request.Content,
			ScheduleTime: request.ScheduleTime,
			Style:        request.Style,
			RepeatTimes:  request.RepeatTimes,
		})
	return db.Error
}
