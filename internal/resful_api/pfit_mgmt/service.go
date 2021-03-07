package pfit_mgmt

import (
	"context"
	"fmt"
	"surf_be/internal/configuration"
	"surf_be/internal/database/redis"
)

type (
	Service struct {
		Config  configuration.Config
		RedisDB redis.Redis
	}
)

func NewService(config configuration.Config, redisDB redis.Redis) Service {
	return Service{
		Config:  config,
		RedisDB: redisDB,
	}
}

func (sv *Service) Login(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	res := LoginResponse{}

	if err := sv.validate(req); err != nil {
		return res, err
	}

	if err := sv.setToRedis(ctx, req); err != nil {
		return res, err
	}

	res.Email = req.Email

	return res, nil
}

func (sv *Service) validate(req LoginRequest) error {
	if err := ValidateStruct(req); err != nil {
		return err
	}
	return nil
}

func (sv *Service) generateKeyInRedis(email string) string {
	return fmt.Sprintf("%v-%v", sv.Config.Environment, email)
}

func (sv *Service) setToRedis(ctx context.Context, req LoginRequest) error {
	if err := sv.RedisDB.Set(ctx, sv.generateKeyInRedis(req.Email), req, 300); err != nil {
		return err
	}
	return nil
}
