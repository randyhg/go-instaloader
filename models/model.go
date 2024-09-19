package models

import (
	"gorm.io/gorm"
	"time"
)

const (
	StatusPending   = 0
	StatusOnProcess = 1
	StatusOk        = 2
	StatusFail      = 3
)

const DefaultStoryLimit = 3

const RedisJobQueueKey = "talent_queue"

type Model struct {
	Id int64 `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
}

type Talent struct {
	Model
	Uuid          string         `gorm:"type:varchar(36)" json:"uuid"`
	Username      string         `gorm:"unique;not null;size:191" json:"username"`
	Url           string         `gorm:"type:text;not null" json:"url"`
	Status        int            `gorm:"default:0" json:"status"`
	StoryImgUrl   string         `gorm:"type:text" json:"story_img_url"`
	StoryImgPath  string         `gorm:"type:text" json:"story_img_path"`
	ProfilePicUrl string         `gorm:"type:text" json:"profile_pic_url"`
	CreatedAt     time.Time      `gorm:"type:timestamp;default:current_timestamp;" json:"-"`
	UpdatedAt     time.Time      `gorm:"type:timestamp;default:current_timestamp ON update current_timestamp;" json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
