package model

import "gorm.io/gorm"

type User struct {
	BaseModel
	Schedule []Schedule `gorm:"foreignKey:UserID"` // 用户拥有的日程
	Username string     `json:"username" gorm:"unique;type:varchar(256);not null"`
	Email    string     `json:"email" gorm:"unique;type:varchar(256);not null"`
	Password string     `json:"password" gorm:"type:varchar(256);not null"`
}

func GetEmailByUserID(db *gorm.DB, userID int) (string, error) {
	var user *User
	db = db.Model(&User{})
	db = db.Where("id = ?", userID)
	result := db.Find(&user)
	return user.Email, result.Error
}
func EmailExist(db *gorm.DB, email string) (bool, error) {
	var count int64
	result := db.Table("users").Where("email = ?", email).Count(&count)
	if count == 0 {
		return false, result.Error
	} else {
		return true, result.Error
	}
}
func UsernameExist(db *gorm.DB, username string) (bool, error) {
	var count int64
	result := db.Table("users").Where("username = ?", username).Count(&count)
	if count == 0 {
		return false, result.Error
	} else {
		return true, result.Error
	}
}
func CreateUser(db *gorm.DB, username, password string, email string) error {
	result := db.Create(&User{
		Username: username,
		Password: password,
		Email:    email,
	})
	return result.Error
}
func GetPasswordByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	return &user, result.Error
}
func GetPasswordByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	return &user, result.Error
}
