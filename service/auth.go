package service

import (
	"context"
	"errors"
	"fmt"
	"rent-car/api/models"
	"rent-car/config"
	"rent-car/pkg"
	"rent-car/pkg/jwt"
	"rent-car/pkg/logger"
	"rent-car/pkg/password"
	"rent-car/pkg/smtp"
	"rent-car/storage"
	"time"
)

type authService struct {
	storage storage.IStorage
	log     logger.ILogger
	redis   storage.IRedisStorage
}

func NewAuthService(storage storage.IStorage, log logger.ILogger, redis storage.IRedisStorage) authService {
	return authService{
		storage: storage,
		log:     log,
		redis:   redis,
	}
}

func (a authService) CustomerLogin(ctx context.Context, loginRequest models.CustomerLoginRequest) (models.CustomerLoginResponse, error) {
	fmt.Println(" loginRequest.Login: ", loginRequest.Login)
	customer, err := a.storage.Customer().GetByLogin(ctx, loginRequest.Login)
	if err != nil {
		a.log.Error("error while getting customer credentials by login", logger.Error(err))
		return models.CustomerLoginResponse{}, err
	}

	if err = password.CompareHashAndPassword(customer.Password, loginRequest.Password); err != nil {
		a.log.Error("error while comparing password", logger.Error(err))
		return models.CustomerLoginResponse{}, err
	}

	m := make(map[interface{}]interface{})

	m["user_id"] = customer.ID
	m["user_role"] = config.CUSTOMER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for customer login", logger.Error(err))
		return models.CustomerLoginResponse{}, err
	}

	return models.CustomerLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a authService) ChangePassword(ctx context.Context, pass models.ChangePassword) (string, error) {
	msg, err := a.storage.Customer().ChangePassword(ctx, pass)
	if err != nil {
		a.log.Error("failed to change customer password", logger.Error(err))
		return "", err
	}

	err = a.redis.Del(ctx, "login:"+pass.Login)
	if err != nil {
		a.log.Error("failed to delete customer data from Redis", logger.Error(err))
	}

	return msg, nil
}

func (a authService) CustomerRegister(ctx context.Context, loginRequest models.CustomerRegisterRequest) error {
	exists, err := a.storage.Customer().CheckEmailExists(ctx, loginRequest.Mail)
	if err != nil {
		a.log.Error("error while checking email existence", logger.Error(err))
		return err
	}
	if exists {
		return errors.New("customer with this email already exists")
	}

	fmt.Println(" loginRequest.Login: ", loginRequest.Mail)
	otpCode := pkg.GenerateOTP()

	msg := fmt.Sprintf("Your OTP code is: %v, for registering RENT_CAR. Don't give it to anyone", otpCode)

	err = a.redis.SetX(ctx, loginRequest.Mail, otpCode, time.Minute*2)
	if err != nil {
		a.log.Error("error while setting otpCode to redis customer register", logger.Error(err))
		return err
	}

	err = smtp.SendMail(loginRequest.Mail, msg)
	if err != nil {
		a.log.Error("error while sending otp code to customer register", logger.Error(err))
		return err
	}

	return nil
}

func (a authService) CustomerRegisterConfirm(ctx context.Context, req models.CustomerRegisterConfirm) (models.CustomerLoginResponse, error) {
	resp := models.CustomerLoginResponse{}

	otp, err := a.redis.Get(ctx, req.Mail)
	if err != nil {
		a.log.Error("error while getting otp code for customer register confirm", logger.Error(err))
		return resp, err
	}
	if req.Otp != otp {
		a.log.Error("incorrect otp code for customer register confirm", logger.Error(err))
		return resp, errors.New("incorrect otp code")
	}

	req.Customer.Email = req.Mail
	id, err := a.storage.Customer().Create(ctx, req.Customer)
	if err != nil {
		a.log.Error("error while creating customer", logger.Error(err))
		return resp, err
	}
	var m = make(map[interface{}]interface{})

	m["user_id"] = id
	m["user_role"] = config.CUSTOMER_ROLE

	accessToken, refreshToken, err := jwt.GenJWT(m)
	if err != nil {
		a.log.Error("error while generating tokens for customer register confirm", logger.Error(err))
		return resp, err
	}
	resp.AccessToken = accessToken
	resp.RefreshToken = refreshToken

	return resp, nil
}
