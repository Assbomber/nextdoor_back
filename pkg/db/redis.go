package db

import (
	"context"
	"fmt"
	"time"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(log *logger.Logger, redisHost, username, password string) *redis.Client {
	log.Info(constants.PENDING + " Connecting to Redis...")
	fmt.Println(redisHost, username, password)
	redisOptions := &redis.Options{
		Addr:     redisHost,
		Username: username,
		Password: password,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client := redis.NewClient(redisOptions)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(constants.FAILURE + " Failed connecting to Redis: " + err.Error())
	}

	log.Info(constants.SUCCESS + " Connected to Redis")
	return client
}
