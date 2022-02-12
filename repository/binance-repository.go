package repository

import (
	"github.com/myomyintko/strategy_robot/model"
	"gorm.io/gorm"
)

type BinanceRepository interface {
	InsertOrder(b model.Order) model.Order
}

type binanceConnection struct {
	connection *gorm.DB
}

func NewBinanceRepository(dbConn *gorm.DB) BinanceRepository {
	return &binanceConnection{
		connection: dbConn,
	}
}

func (db *binanceConnection) InsertOrder(order model.Order) model.Order {
	db.connection.Save(&order)
	db.connection.Preload("Robot").Find(&order)
	return order
}
