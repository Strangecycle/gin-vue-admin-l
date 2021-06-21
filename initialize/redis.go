package initialize

import (
	"gin-vue-admin-l/global"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func Redis() {
	redisConf := global.GVA_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})
	pong, err := client.Ping().Result()
	if err != nil {
		global.GVA_LOG.Error("redis connect ping failed, err:", zap.Any("err", err))
	} else {
		global.GVA_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.GVA_REDIS = client
	}
}
