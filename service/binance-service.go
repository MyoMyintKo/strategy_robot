package service

import (
	"log"

	"github.com/mashingan/smapping"
	"github.com/myomyintko/strategy_robot/dto"
	"github.com/myomyintko/strategy_robot/model"
	"github.com/myomyintko/strategy_robot/repository"
)

type BinanceService interface {
	Insert(b dto.CreateOrderDTO) model.Order
}

type binanceService struct {
	binanceRepository repository.BinanceRepository
}

func NewBinanceService(binRepo repository.BinanceRepository) BinanceService {
	return &binanceService{
		binanceRepository: binRepo,
	}
}

func (service *binanceService) Insert(o dto.CreateOrderDTO) model.Order {
	order := model.Order{}
	err := smapping.FillStruct(&order, smapping.MapFields(&o))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.binanceRepository.InsertOrder(order)
	return res
}
