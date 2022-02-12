package dto

type RobotUpdateDTO struct {
	ID     uint64 `json:"id" form:"id"`
	Symbol   string `json:"symbol" form:"symbol" binding:"required"`
	UserID uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}

type RobotCreateDTO struct {
	Symbol   string `json:"symbol" form:"symbol" binding:"required"`
	UserID uint64 `json:"user_id,omitempty"  form:"user_id,omitempty"`
}
