package pfit_mgmt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"surf_be/internal/app/bot"
	"surf_be/internal/app/utils"
	"surf_be/internal/configuration"
	"surf_be/internal/database/redis"
	"surf_be/internal/resful_api/binance"
	binanceWS "surf_be/internal/websocket/binance"
	"time"
)

const (
	LoginTTLDefaut = 10080
)

type (
	Service struct {
		Config      configuration.Config
		RedisDB     *redis.Redis
		Repo        Repository
		BinanceSv   *binance.Service
		BinanceWSHL *binanceWS.HandlerImpl
		Bots        []string
	}
)

func NewService(config configuration.Config, redisDB *redis.Redis, repo Repository, binanceSv *binance.Service, binanceWSHL *binanceWS.HandlerImpl) Service {
	return Service{
		Config:      config,
		RedisDB:     redisDB,
		Repo:        repo,
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

	res.EmailWithEnv = fmt.Sprintf("%v-%v", sv.Config.Environment, req.Email)

	return res, nil
}

func (sv *Service) validate(req interface{}) error {
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
		if err := sv.RedisDB.Set(ctx, sv.generateKeyInRedis(req.Email), req, LoginTTLDefaut); err != nil {
			return err
		}
	} else {
		if err := sv.RedisDB.Delete(ctx, sv.generateKeyInRedis(req.Email)); err != nil {
			return err
		}
		if err := sv.RedisDB.Set(ctx, sv.generateKeyInRedis(req.Email), req, LoginTTLDefaut); err != nil {
			return err
		}
	}
	return nil
}

func (sv *Service) StartBot(ctx context.Context, req StartBotRequest) (StartBotResponse, error) {
	res := StartBotResponse{}

	if err := sv.validate(req); err != nil {
		return res, err
	}

	email := ctx.Value(ContextEmailKey).(string)

	pair := fmt.Sprintf("%v%v", req.Access, req.Quote)

	botDB, err := sv.Repo.GetBotByEmailAndPair(ctx, email, pair)
	if err != nil {
		return res, err
	}

	if botDB != nil {
		return res, errors.New("bot is exited")
	}

	newBot, err := sv.newBot(req.Access, pair)
	if err != nil {
		return res, err
	}
	newBot.Email = email
	if err := sv.Repo.InsertBot(ctx, *newBot); err != nil {
		return res, err
	}

	sv.Bots = append(sv.Bots, newBot.ID)

	res.BotID = newBot.ID
	return res, nil
}

func (sv *Service) newBot(access, pair string) (*bot.Bot, error) {
	rspData, err := sv.BinanceSv.GetAggTrades(pair, "1")
	if err != nil {
		return nil, err
	}

	currentPrice, err := strconv.ParseFloat(rspData.Price, 32)
	if err != nil {
		return nil, err
	}

	TempBot := bot.Bot{
		ID:            sv.generateBotID(),
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

	sv.BinanceWSHL.PushBot(&TempBot)
	return &TempBot, nil
}

func (sv *Service) GetGetBotStatus(ctx context.Context, req GetBotStatusRequest) (GetBotStatusResponse, error) {
	res := GetBotStatusResponse{}

	if err := sv.validate(req); err != nil {
		return res, err
	}

	email := ctx.Value(ContextEmailKey).(string)

	pair := fmt.Sprintf("%v%v", req.Access, req.Quote)

	botDB, err := sv.Repo.GetBotByEmailAndPair(ctx, email, pair)
	if err != nil {
		return res, err
	}

	if botDB == nil {
		return res, errors.New("bot is not existed")
	}

	botRunning := sv.BinanceWSHL.GetBot(botDB.ID)

	res.CurrentPrice = botRunning.GetCurrentData()
	return res, nil
}

func (sv *Service) generateBotID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
