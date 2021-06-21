package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gorm.io/gorm"
	"time"
)

// 将 jwt 拉入黑名单作废
func JsonInBlacklist(jwtBlack model.JwtBlacklist) (err error) {
	err = global.GVA_DB.Create(&jwtBlack).Error
	return err
}

// token 是否在黑名单中
func IsInBlacklist(jwt string) bool {
	err := global.GVA_DB.Where("jwt = ?", jwt).First(&model.JwtBlacklist{}).Error
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// 从 redis 中获取 jwt
func GetRedisJWT(username string) (err error, redisJWT string) {
	result, err := global.GVA_REDIS.Get(username).Result()
	return err, result
}

// jwt 存入 redis 并设置过期时间
func SetRedisJWT(jwt string, username string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.GVA_CONFIG.JWT.ExpiresTime) * time.Second
	err = global.GVA_REDIS.Set(username, jwt, timer).Err()
	return err
}
