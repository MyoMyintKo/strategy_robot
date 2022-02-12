package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/myomyintko/strategy_robot/dto"
	"github.com/myomyintko/strategy_robot/model"
	"github.com/myomyintko/strategy_robot/repository"
)

//UserService is a contract.....
type UserService interface {
	Update(user dto.UserUpdateDTO) model.User
	Profile(userID string) model.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) model.User {
	userToUpdate := model.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed map %v:", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userID string) model.User {
	return service.userRepository.ProfileUser(userID)
}
