package model

import (
	"time"
)

type BinanceAPI struct {
	ID        uint64 `gorm:"primary_key:auto_increment" json:"id"`
	APIKey    string `gorm:"unique,type:varchar(255)" json:"api"`
	SecretKey string `gorm:"unique,type:varchar(255)" json:"secret"`
	StreamKey string `gorm:"unique,type:varchar(255)" json:"stream"`
	UserID    uint64 `gorm:"not null" json:"-"`
	User      User   `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	BoundAt   time.Time
	StreamedAt time.Time
}

type Order struct {
	ID               uint64  `gorm:"primary_key:autoincrement" json:"id"`
	OrderId int64 `json:"order_id"`
	ClientOrderId string  `json:"client_order_id"`
	RobotID           uint64  `gorm:"not null" json:"-"`
	Robot             Robot    `gorm:"foreignKey:RobotID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"robot"`
	OrderedAt time.Time
}
