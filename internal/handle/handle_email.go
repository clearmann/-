package handle

import (
	g "backend/internal/global"
	"backend/internal/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

var (
	ScheduleMap sync.Map
	taskMutex   sync.Mutex // 用于保护发邮件操作
	DB          *gorm.DB
)

func SendEmail(dialer *gomail.Dialer, schedule *model.Schedule) {
	log.Println("准备发邮件...")
	m := gomail.NewMessage()
	m.SetHeader("From", g.Conf.Email.EmailFrom)
	m.SetHeader("To", schedule.Email)
	m.SetHeader("Subject", schedule.Title)
	m.SetBody("text/html", schedule.Content)
	// 设置SMTP服务器信息
	// 使用DialAndSend方法发送邮件
	if err := dialer.DialAndSend(m); err != nil {
		panic(err.Error())
	}
	log.Println("Email sent successfully!")
	handleScheduleType(schedule)
}
func handleScheduleType(schedule *model.Schedule) {
	switch schedule.Style {
	case 1: //1 代表此日程 只执行一次
		ScheduleMap.Delete(schedule.ID) //删除日程
		err := model.DelSchedule(DB, schedule.UserID, schedule.ID)
		log.Printf("scheduleID: %d , UserID: %d", schedule.ID, schedule.UserID)
		if err != nil {
			log.Println("删除操作失败...")
			return
		}
	case 2: //2 代表此日程每天都要执行 需在本次执行完之后将时间增加24小时
		schedule.ScheduleTime = schedule.ScheduleTime.Add(24 * time.Hour)
		ScheduleMap.Store(schedule.ID, schedule)
		err := model.UpdateSchedule(DB, schedule)
		if err != nil {
			log.Println("更新日程失败...")
			return
		}
	case 3: //3 代表此日程周一到周五每天执行  如果是周五，则跳到周一；否则增加24小时
		if schedule.ScheduleTime.Weekday() == time.Friday {
			schedule.ScheduleTime = schedule.ScheduleTime.Add(3 * time.Hour)
		} else {
			schedule.ScheduleTime = schedule.ScheduleTime.Add(24 * time.Hour)
		}
		ScheduleMap.Store(schedule.ID, schedule)
		err := model.UpdateSchedule(DB, schedule)
		if err != nil {
			log.Println("更新日程失败...")
			return
		}
	default:
		log.Println("未记录的schedule style")
		return
	}
}
func InitEmailScheduler(c *gin.Context) {
	DB = GetDB(c)
	log.Println("InitEmailScheduler...")
	dialer := gomail.NewDialer(g.Conf.Email.SMTPHost, g.Conf.Email.SMTPPort, g.Conf.Email.EmailFrom, g.Conf.Email.EmailPassword)
	emailWorkerNumber := g.Conf.Email.EmailWorkerNumber
	ticker := time.NewTicker(10 * time.Second) // 设置为每十秒触发一次  可根据具体业务进行调整
	//defer ticker.Stop()
	for i := 0; i < emailWorkerNumber; i++ {
		go func(id int) {
			for range ticker.C {
				now := time.Now()
				ScheduleMap.Range(func(k, value any) bool {
					log.Println("循环sync.map")
					schedule := value.(*model.Schedule)
					//当前时间或晚一段时间发，可根据具体场景具体定义
					if now.Equal(schedule.ScheduleTime) || now.After(schedule.ScheduleTime) {
						taskMutex.Lock()
						SendEmail(dialer, schedule) // 发送邮件
						taskMutex.Unlock()
					}
					return true
				})
			}
		}(i)
	}
}
