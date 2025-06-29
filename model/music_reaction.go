package model

import "time"

type MusicReaction struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"not null;index;uniqueIndex:idx_user_music" json:"user_id"`
	MusicUUID string    `gorm:"type:varchar(64);not null;index;uniqueIndex:idx_user_music" json:"uuid"`
	Action    int64     `gorm:"not null" json:"action"` // 行为类型：1=点赞，0=未点赞，其他值可扩展为收藏这些功能
	CreatedAt time.Time `json:"created_at"`

	User      User      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	MusicFile MusicFile `gorm:"foreignKey:MusicUUID;references:UUID" json:"-"`
}
