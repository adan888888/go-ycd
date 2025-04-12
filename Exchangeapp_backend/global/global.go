package global

import (
	"exchangeapp/models"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Db        *gorm.DB
	RedisDB   *redis.Client
	AppConfig *models.Config
	//CurrentTempIndex int
)
