package model

import "time"

type Robot struct {
	ID        uint64   `gorm:"primary_key:auto_increment" json:"id"`
	Symbol    string   `gorm:"type:varchar(255)" json:"symbol"`
	UserID    uint64   `gorm:"not null" json:"-"`
	User      User     `gorm:"foreignKey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user"`
	Orders    *[]Order `json:"orders,omitempty"`
	CreatedAt time.Time
}
