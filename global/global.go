package global

import (
	"gin-vue-admin-l/config"
	"gin-vue-admin-l/utils/timer"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GVA_CONFIG config.Server
	GVA_VP     *viper.Viper
	GVA_LOG    *zap.Logger
	GVA_DB     *gorm.DB
	GVA_TIMER  timer.Timer = timer.NewTimerTask()
	GVA_REDIS  *redis.Client
)
