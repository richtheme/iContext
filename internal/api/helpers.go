package api

import (
	"iContext/repository/postgres"
	"iContext/repository/redis"

	"github.com/sirupsen/logrus"
)

func (a *API) configureLoggerField() error {
	logLevel, err := logrus.ParseLevel(a.config.LoggerLevel)
	if err != nil {
		return err
	}
	a.logger.SetLevel(logLevel)
	return nil
}

func (a *API) configureRouterField() {
	a.router.POST("/redis/incr", a.redisIncr)
	a.router.POST("/sign/hmacsha512", a.cryptHmac)
	a.router.POST("/postgres/users", a.SaveUser)
}

func (a *API) configureStorageField() error {
	storage := postgres.New(a.config.Storage)
	if err := storage.Open(); err != nil {
		return err
	}

	a.storage = storage
	return nil
}

func (a *API) configureRedisField() error {
	client, err := redis.RedisNew(a.redisAddr)
	if err != nil {
		return err
	}

	a.client = client
	return nil
}
