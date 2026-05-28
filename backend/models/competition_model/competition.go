package competition_model

import "time"

type Competition struct {
	ID     int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"column:name;type:varchar(100);not null;comment:竞赛名称"`
	Status int8   `json:"status" gorm:"column:status;type:tinyint;not null;default:1;comment:1未开始 2进行中 3已结束"`

	Visibility   bool   `json:"visibility" gorm:"column:visibility;type:tinyint(1);not null;comment:是否公开 0非公开 1公开"`
	PasswordHash string `json:"-" gorm:"column:password_hash;type:varchar(255);comment:非公开竞赛密码哈希"`

	PlayerCount int `json:"player_count" gorm:"column:player_count;not null;default:0;comment:总参赛人数/报名人数"`

	StartTime time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime   time.Time `json:"end_time" gorm:"column:end_time"`

	CreatorID   int64  `json:"creator_id" gorm:"column:creator_id;not null"`
	CreatorName string `json:"creator_name" gorm:"column:creator_name;type:varchar(64);not null"`

	Announcement string `json:"announcement" gorm:"column:announcement;type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

type CompetitionsResp struct {
	Competitions []*Competition `json:"competitions"`
}

type GetCompetitionReq struct {
	CompetitionID int64 `json:"competition_id" uri:"competition_id" binding:"required"`
}

type CompetitionResp struct {
	Competitions Competition `json:"competition"`
}
