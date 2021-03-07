package pfit_mgmt

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"surf_be/internal/app/bot"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
	"surf_be/internal/database/redis"
	"surf_be/internal/resful_api/binance"
	binanceWS "surf_be/internal/websocket/binance"
	"time"
)

type (
	Service struct {
		Config      configuration.Config
		RedisDB     redis.Redis
		BinanceSv   *binance.Service
		BinanceWSHL *binanceWS.HandlerImpl
		IsStarted   bool
		Bot         *bot.Bot
	}
)

func NewService(config configuration.Config, redisDB redis.Redis, binanceSv *binance.Service, binanceWSHL *binanceWS.HandlerImpl) Service {
	return Service{
		Config:      config,
		RedisDB:     redisDB,
		BinanceSv:   binanceSv,
		BinanceWSHL: binanceWSHL,
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

	if !sv.IsStarted {
		sv.Bot = sv.startBot()
		sv.IsStarted = true
	}
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
	value, err := sv.RedisDB.Get(ctx, sv.generateKeyInRedis(req.Email))
	if err != nil {
		return err
	}
	if value == nil {
		if err := sv.RedisDB.Set(ctx, sv.generateKeyInRedis(req.Email), req, 300); err != nil {
			return err
		}
	} else {
		if err := sv.RedisDB.Delete(ctx, sv.generateKeyInRedis(req.Email)); err != nil {
			return err
		}
		if err := sv.RedisDB.Set(ctx, sv.generateKeyInRedis(req.Email), req, 300); err != nil {
			return err
		}
	}
	return nil
}

func (sv *Service) startBot() *bot.Bot {
	access := "BTC"
	excess := "USDT"
	pair := fmt.Sprintf("%v%v", access, excess)

	rspData, err := sv.BinanceSv.GetAggTrades(pair, "1")
	if err != nil {
		log.Fatal(err)
	}

	currentPrice, err := strconv.ParseFloat(rspData.Price, 32)
	if err != nil {
		log.Fatalf("error parse float: %v", err)
	}

	BTCBot := bot.Bot{
		ID:            1,
		StartTime:     time.Now(),
		Duration:      2 * time.Hour,
		Pair:          rspData.Symbol,
		Access:        access,
		BuyInPrice:    currentPrice,
		BuyInQuantity: 1,
		CurrentPrice:  currentPrice,
		Quantity:      1,
		StopChannel:   nil,
		Type:          utils.AggTradeStreamType,
		PercentBuy:    0.01,
		Budget:        0,
	}

	sv.BinanceWSHL.PushBot(&BTCBot)
	return &BTCBot
}

func (sv *Service) GetGetBotStatus(ctx context.Context) (string, error) {
	return sv.Bot.GetCurrentData(), nil
}
