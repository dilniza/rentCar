package service

import (
	"context"
	"encoding/json"
	"rent-car/api/models"
	"rent-car/pkg/logger"
	"rent-car/storage"
)

type customerService struct {
	storage storage.IStorage
	logger  logger.ILogger
	redis   storage.IRedisStorage
}

func NewCustomerService(storage storage.IStorage, logger logger.ILogger, redis storage.IRedisStorage) customerService {
	return customerService{
		storage: storage,
		logger:  logger,
		redis:   redis,
	}
}

func (s customerService) Create(ctx context.Context, customer models.CreateCustomer) (string, error) {
	pKey, err := s.storage.Customer().Create(ctx, customer)
	if err != nil {
		s.logger.Error("failed to create customer", logger.Error(err))
		return "", err
	}

	return pKey, nil
}

func (s customerService) Update(ctx context.Context, customer models.UpdateCustomer, id string) (string, error) {
	id, err := s.storage.Customer().Update(ctx, customer, id)
	if err != nil {
		s.logger.Error("failed to update customer", logger.Error(err))
		return "", err
	}
	return id, nil
}

func (s customerService) GetByID(ctx context.Context, id string) (models.Customer, error) {
	var customer models.Customer

	customerData, err := s.redis.Get(ctx, "customer_id:"+id)
	if err == nil {
		err := json.Unmarshal([]byte(customerData.(string)), &customer)
		if err != nil {
			s.logger.Error("failed to unmarshal customer data from Redis", logger.Error(err))
		} else {
			return customer, nil
		}
	}

	customer, err = s.storage.Customer().GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get customer by ID", logger.Error(err))
		return models.Customer{}, err
	}
	return customer, nil
}

func (s customerService) GetAll(ctx context.Context, req models.GetAllCustomersRequest) (models.GetAllCustomersResponse, error) {

	customers, err := s.storage.Customer().GetAll(ctx, req)
	if err != nil {
		s.logger.Error("failed to get all customers", logger.Error(err))
		return models.GetAllCustomersResponse{}, err
	}

	return customers, nil
}

func (s customerService) GetCustomerCars(ctx context.Context, name string, id string, boolean bool) (models.GetCustomerCarsResponse, error) {

	customerCars, err := s.storage.Customer().GetCustomerCars(ctx, name, id, boolean)
	if err != nil {
		s.logger.Error("failed to get customer cars", logger.Error(err))
		return models.GetCustomerCarsResponse{}, err
	}

	return customerCars, nil
}

func (s customerService) Delete(ctx context.Context, id string) error {

	err := s.storage.Customer().Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete customer", logger.Error(err))
		return err
	}

	err = s.redis.Del(ctx, "customer_id:"+id)
	if err != nil {
		s.logger.Error("failed to delete customer data from Redis", logger.Error(err))
	}

	return nil
}
