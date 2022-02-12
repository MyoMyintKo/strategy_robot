package model

import "time"

type User struct {
	ID       uint64   `gorm:"primary_key:auto_increment" json:"id"`
	Name     string   `gorm:"type:varchar(255)" json:"name"`
	Email    string   `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string   `gorm:"->;<-;not null" json:"-"`
	Token    string   `gorm:"-" json:"token,omitempty"`
	Robots   []*Robot `json:"robots,omitempty"`
	Keys   []*BinanceAPI `json:"keys,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
