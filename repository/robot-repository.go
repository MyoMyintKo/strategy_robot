package repository

import (
	"github.com/myomyintko/strategy_robot/model"
	"gorm.io/gorm"
)

type RobotRepository interface {
	InsertRobot(b model.Robot) model.Robot
	UpdateRobot(b model.Robot) model.Robot
	DeleteRobot(b model.Robot)
	AllRobot() []model.Robot
	FindRobotByID(robotID uint64) model.Robot
	FindRobotByUserID(robotID uint64) model.Robot
}

type robotConnection struct {
	connection *gorm.DB
}

func NewRobotRepository(dbConn *gorm.DB) RobotRepository {
	return &robotConnection{
		connection: dbConn,
	}
}

func (db *robotConnection) InsertRobot(robot model.Robot) model.Robot {
	db.connection.Save(&robot)
	db.connection.Preload("User").Find(&robot)
	return robot
}

func (db *robotConnection) UpdateRobot(robot model.Robot) model.Robot {
	db.connection.Save(&robot)
	db.connection.Preload("User").Find(&robot)
	return robot
}

func (db *robotConnection) DeleteRobot(robot model.Robot) {
	db.connection.Delete(&robot)
}

func (db *robotConnection) FindRobotByID(robotID uint64) model.Robot {
	var robot model.Robot
	db.connection.Preload("User").Find(&robot, robotID)
	return robot
}

func (db *robotConnection) FindRobotByUserID(robotID uint64) model.Robot {
	var robot model.Robot
	db.connection.Preload("User").Where("user_id = ?",robotID).Find(&robot)
	return robot
}

func (db *robotConnection) AllRobot() []model.Robot {
	var robots []model.Robot
	db.connection.Preload("User").Find(&robots)
	return robots
}
