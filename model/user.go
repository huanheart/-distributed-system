package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(50)" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`   // 唯一索引
	Username  string         `gorm:"type:varchar(50);uniqueIndex" json:"username"` // 唯一索引
	Password  string         `gorm:"type:varchar(255)" json:"-"`                   // 不返回给前端
	CreatedAt time.Time      `json:"created_at"`                                   // 自动时间戳
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 支持软删除
}
