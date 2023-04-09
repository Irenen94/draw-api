package model

import "time"

type AlgoLog struct {
	ID        int64  `gorm:"column:id" db:"id" json:"id" form:"id"`
	RequestId string `gorm:"column:request_id" db:"request_id" json:"request_id" form:"request_id"`
	Prompt    string `gorm:"column:prompt" db:"prompt" json:"prompt" form:"prompt"`
	TaskId    string `gorm:"column:task_id" db:"task_id" json:"task_id" form:"task_id"`             //  证件类型
	DdimSteps int    `gorm:"column:ddim_steps" db:"ddim_steps" json:"ddim_steps" form:"ddim_steps"` //  证件类型
	Artists   string `gorm:"column:artists"  db:"artists" json:"artists" form:"artists"`
	Styles    string `gorm:"column:styles" db:"styles" json:"styles" form:"styles" `
	SizeType  int    `gorm:"column:size_type" db:"size_type" json:"size_type" form:"size_type"`
	H            int       `gorm:"column:h" db:"h" json:"h" form:"h"`
	W            int       `gorm:"column:w" db:"w" json:"w" form:"w"`
	NSamples     int       `gorm:"column:n_samples" db:"n_samples" json:"n_samples" form:"n_samples"`
	Negative     string    `gorm:"column:negative" db:"negative" json:"negative" form:"negative"`
	ImageType    string    `gorm:"column:image_type" db:"image_type" json:"image_type" form:"image_type"`
	SmartMode    string    `gorm:"column:smart_mode" db:"smart_mode" json:"smart_mode" form:"smart_mode"`
	Seed         int       `gorm:"column:seed" db:"seed" json:"seed" form:"seed"`
	InitImage    string    `gorm:"column:init_image" db:"init_image" json:"init_image" form:"init_image"`
	CreateImages string    `gorm:"column:create_images" db:"create_images" json:"create_images" form:"create_images"`
	SnapTime     time.Time `gorm:"column:snap_time" db:"snap_time" json:"snap_time" form:"snap_time"` //  调用时间，毫秒时间戳
	Code         int       `gorm:"column:code" db:"code" json:"code" form:"code"`                     //  接口返回状态码
	Message      string    `gorm:"column:message" db:"message" json:"message" form:"message"`         //  接口返回状态信息
	Details      string    `gorm:"column:details" db:"details" json:"details" form:"details"`
	TotalTime    int       `gorm:"column:total_time" db:"total_time" json:"total_time" form:"total_time"`
}
