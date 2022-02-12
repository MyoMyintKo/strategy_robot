package dto

import "time"

type APIUpdateDTO struct {
	ID        uint64 `json:"id" form:"id"`
	APIKey    string `json:"api" form:"api" binding:"required"`
	SecretKey string `json:"secret" form:"secret" binding:"required"`
	UserID    uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}

type BindStreamDTO struct {
	ID        uint64 `json:"id" form:"id"`
	UserID    uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
	StreamKey string `json:"stream" form:"stream" binding:"required"`
	StreamedAt time.Time
}

type APICreateDTO struct {
	APIKey    string `json:"api" form:"api" binding:"required"`
	SecretKey string `json:"secret" form:"secret" binding:"required"`
	UserID    uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
	BoundAt time.Time
}

type CreateOrderDTO struct {
	OrderId       int64  `json:"order_id"`
	ClientOrderId string `json:"client_order_id"`
	Price         string `json:"price" form:"price" binding:"required"`
	Quantity      string `json:"quantity" form:"quantity" binding:"required"`
	RobotID       uint64 `json:"robot_id,omitempty" form:"robot_id,omitempty"`
}
