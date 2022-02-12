package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/myomyintko/strategy_robot/dto"
	"github.com/myomyintko/strategy_robot/model"
	"github.com/myomyintko/strategy_robot/repository"
)

type RobotService interface {
	Insert(b dto.RobotCreateDTO) model.Robot
	Update(b dto.RobotUpdateDTO) model.Robot
	Delete(b model.Robot)
	All() []model.Robot
	FindByID(robotID uint64) model.Robot
	FindByUserID(userID uint64) model.Robot
	IsUserExistRobot(symbol string, userID uint64) bool
	IsAllowedToEdit(userID, robotID uint64) bool
}

type robotService struct {
	robotRepository repository.RobotRepository
}

func NewRobotService(robotRepo repository.RobotRepository) RobotService {
	return &robotService{
		robotRepository: robotRepo,
	}
}

func (service *robotService) Insert(b dto.RobotCreateDTO) model.Robot {
	robot := model.Robot{}
	err := smapping.FillStruct(&robot, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.robotRepository.InsertRobot(robot)
	return res
}

func (service *robotService) Update(b dto.RobotUpdateDTO) model.Robot {
	robot := model.Robot{}
	err := smapping.FillStruct(&robot, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.robotRepository.UpdateRobot(robot)
	return res
}

func (service *robotService) Delete(b model.Robot) {
	service.robotRepository.DeleteRobot(b)
}

func (service *robotService) All() []model.Robot {
	return service.robotRepository.AllRobot()
}

func (service *robotService) FindByID(robotID uint64) model.Robot {
	return service.robotRepository.FindRobotByID(robotID)
}

func (service *robotService) FindByUserID(userID uint64) model.Robot {
	return service.robotRepository.FindRobotByUserID(userID)
}

func (service *robotService) IsUserExistRobot(symbol string, userID uint64) bool {
	robot := service.robotRepository.FindRobotByUserID(userID)
	sym := robot.Symbol
	return symbol == sym
}

func (service *robotService) IsAllowedToEdit(userID, robotID uint64) bool {
	robot := service.robotRepository.FindRobotByID(robotID)
	id := robot.UserID
	return userID == id
}
