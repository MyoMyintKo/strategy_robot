package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/myomyintko/strategy_robot/dto"
	"github.com/myomyintko/strategy_robot/model"
	"github.com/myomyintko/strategy_robot/repository"
)

type APIService interface {
	Insert(b dto.APICreateDTO) model.BinanceAPI
	Update(b dto.APIUpdateDTO) model.BinanceAPI
	BindStream(b dto.BindStreamDTO) model.BinanceAPI
	Delete(b model.BinanceAPI)
	All() []model.BinanceAPI
	FindByID(apiID uint64) model.BinanceAPI
	FindByUserID(userID uint64) model.BinanceAPI
	IsAllowedToEdit(userID , apiID uint64) bool
}

type apiService struct {
	apiRepository repository.APIRepository
}

func NewAPIService(apiRepo repository.APIRepository) APIService {
	return &apiService{
		apiRepository: apiRepo,
	}
}

func (service *apiService) Insert(b dto.APICreateDTO) model.BinanceAPI {
	api := model.BinanceAPI{}
	err := smapping.FillStruct(&api, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.apiRepository.InsertAPI(api)
	return res
}

func (service *apiService) Update(b dto.APIUpdateDTO) model.BinanceAPI {
	key := model.BinanceAPI{}
	err := smapping.FillStruct(&key, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.apiRepository.UpdateAPI(key)
	return res
}

func (service *apiService) BindStream(b dto.BindStreamDTO) model.BinanceAPI{
	key := model.BinanceAPI{}
	err := smapping.FillStruct(&key, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.apiRepository.BindStream(key)
	return res
}

func (service *apiService) Delete(b model.BinanceAPI) {
	service.apiRepository.DeleteAPI(b)
}

func (service *apiService) All() []model.BinanceAPI {
	return service.apiRepository.AllAPI()
}

func (service *apiService) FindByID(apiID uint64) model.BinanceAPI {
	return service.apiRepository.FindAPIByID(apiID)
}

func (service *apiService) FindByUserID(userID uint64) model.BinanceAPI {
	return service.apiRepository.FindAPIByUserID(userID)
}

func (service *apiService) IsAllowedToEdit(userID , apiID uint64) bool {
	b := service.apiRepository.FindAPIByID(apiID)
	id := b.UserID
	return userID == id
}
