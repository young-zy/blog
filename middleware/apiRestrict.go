package middleware

import (
	"blog/conf"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	restrict "github.com/young-zy/gin-api-restriction"
)

// ApiRestrict is a restriction middleware
var ApiRestrict *restrict.RestrictionMiddleWare

func init() {

	config := conf.Config

	restrictConf := config.Restrict
	redisConf := config.Redis

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Addr, redisConf.Port),
		Username: redisConf.Username,
		Password: redisConf.Password,
		DB:       redisConf.DB,
	}

	redisClient := redis.NewClient(options)

	ApiRestrict = restrict.NewDefaultRestrictionMiddleWare(&restrict.RestrictionConfig{
		Log:              false,
		RestrictionCount: restrictConf.RestrictionCount,
		RestrictionTime:  time.Duration(restrictConf.RestrictionTime) * time.Second,
	}, redisClient)
}
