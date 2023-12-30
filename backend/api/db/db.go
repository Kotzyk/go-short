package db

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func readDefaultRateLimitFromEnv() uint {
	limit, err := strconv.Atoi(os.Getenv("DEFAULT_RATELIMIT"))
	if err != nil {
		panic(err)
	}
	return uint(limit)
}
func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

var Ctx = context.Background()

func CreateClient(dbInt int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbInt,
	})

	return rdb
}

func RateLimitMiddleware(rdb *redis.Client) gin.HandlerFunc {

	limitedStore := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: rdb,
		Rate:        time.Second,
		Limit:       readDefaultRateLimitFromEnv(),
	})

	mw := ratelimit.RateLimiter(limitedStore, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
	return mw
}
