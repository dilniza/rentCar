package service

import (
	"rent-car/pkg/logger"
	"rent-car/storage"
)

type IServiceManager interface {
	Car() carService
	Customer() customerService
	Order() orderService
	Auth() authService
}

type Service struct {
	carService      carService
	customerService customerService
	orderService    orderService
	auth            authService

	logger logger.ILogger
}

func New(storage storage.IStorage, log logger.ILogger, redis storage.IRedisStorage) Service {
	return Service{
		carService:      NewCarService(storage, log),
		customerService: NewCustomerService(storage, log, redis),
		orderService:    NewOrderService(storage, log),
		auth:            NewAuthService(storage, log, redis),
		logger:          log,
	}
}

func (s Service) Car() carService {
	return s.carService
}

func (s Service) Customer() customerService {
	return s.customerService
}

func (s Service) Order() orderService {
	return s.orderService
}

func (s Service) Auth() authService {
	return s.auth
}
