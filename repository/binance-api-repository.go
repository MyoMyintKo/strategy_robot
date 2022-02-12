package repository

import (
	"fmt"
	"github.com/myomyintko/strategy_robot/model"
	"gorm.io/gorm"
)

type APIRepository interface {
	InsertAPI(b model.BinanceAPI) model.BinanceAPI
	UpdateAPI(b model.BinanceAPI) model.BinanceAPI
	BindStream(b model.BinanceAPI) model.BinanceAPI
	DeleteAPI(b model.BinanceAPI)
	AllAPI() []model.BinanceAPI
	FindAPIByID(id uint64) model.BinanceAPI
	FindAPIByUserID(userID uint64) model.BinanceAPI
}

type apiConnection struct {
	connection *gorm.DB
}

func NewAPIRepository(dbConn *gorm.DB) APIRepository {
	return &apiConnection{
		connection: dbConn,
	}
}

func (db *apiConnection) InsertAPI(key model.BinanceAPI) model.BinanceAPI {
	db.connection.Save(&key)
	db.connection.Preload("User").Find(&key)
	return key
}

func (db *apiConnection) UpdateAPI(key model.BinanceAPI) model.BinanceAPI {
	db.connection.Save(&key)
	db.connection.Preload("User").Find(&key)
	return key
}

func (db *apiConnection) BindStream(key model.BinanceAPI) model.BinanceAPI {
	if err := db.connection.Model(&key).Where("user_id = ?", &key.UserID).Updates(&key).Error; err != nil{
		fmt.Println(err)
	}
	db.connection.Preload("User").Find(&key)
	return key
}

func (db *apiConnection) DeleteAPI(key model.BinanceAPI) {
	db.connection.Delete(&key)
}

func (db *apiConnection) FindAPIByID(id uint64) model.BinanceAPI {
	var key model.BinanceAPI
	db.connection.Preload("User").Find(&key, id)
	return key
}

func (db *apiConnection) FindAPIByUserID(userID uint64) model.BinanceAPI {
	var key model.BinanceAPI
	db.connection.Find(&key, "user_id =?",userID)
	return key
}

func (db *apiConnection) AllAPI() []model.BinanceAPI {
	var keys []model.BinanceAPI
	db.connection.Preload("User").Find(&keys)
	return keys
}
