package api

import (
	"iContext/repository/postgres"
	"iContext/repository/redis"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type API struct {
	config    *Config
	logger    *logrus.Logger
	router    *gin.Engine
	storage   *postgres.Storage
	redisAddr string
	client    *redis.RedisStorage
}

func New(config *Config, redisAddr string) *API {
	return &API{
		config:    config,
		logger:    logrus.New(),
		router:    gin.Default(),
		redisAddr: redisAddr,
	}
}

func (api *API) Start() error {
	if err := api.configureLoggerField(); err != nil {
		return err
	}
	api.logger.Info("starting api server at port", api.config.BindAddr)

	api.configureRouterField()

	if err := api.configureStorageField(); err != nil {
		return err
	}

	createTable := `CREATE TABLE users
	(
		id serial not null unique,
		name varchar(255) not null,
		age int not null
	);`
	_ = api.storage.Exec(createTable)

	if err := api.configureRedisField(); err != nil {
		return err
	}

	return api.router.Run(api.config.BindAddr)
}
