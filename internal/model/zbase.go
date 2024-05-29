package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MakeMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&User{},
		&Schedule{},
	)
	return err
}

// 分页
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// 数据列表
func List[T any](db *gorm.DB, data T, slt, order, query string, args ...any) (T, error) {
	db = db.Model(data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}
	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}

// 统计数量
func Count[T any](db *gorm.DB, data *T, where ...any) (int, error) {
	var total int64
	db = db.Model(data)
	if len(where) > 0 {
		db = db.Where(where[0], where[1:]...)
	}
	result := db.Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(total), nil
}
