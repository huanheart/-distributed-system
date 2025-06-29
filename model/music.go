package model

import (
	"gorm.io/gorm"
	"time"
)

// gorm 是默认驼峰型转换成蛇形的
// 即大小写交界处才会转化成下划线_
// 固然UUID依旧是uuid
type MusicFile struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	UUID      string         `gorm:"type:varchar(64);uniqueIndex" json:"uuid"` // 物理文件唯一标识
	UserID    int64          `gorm:"index" json:"user_id"`                     // 上传用户 ID
	MusicName string         `gorm:"type:varchar(100)" json:"name"`            // 原始音乐名（展示用）
	FilePath  string         `gorm:"type:varchar(255)" json:"file_path"`       // 存储路径或 URL
	LikeCount int64          `json:"like_count"`
	FileSize  int64          `json:"file_size"` // 文件大小（字节）
	Duration  float64        `json:"duration"`  // 播放时长（秒）
	IsUpload  int64          `gorm:"column:is_upload;" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
