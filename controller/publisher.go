package controller

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type RequestPublisher struct {
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

var redisc = func(host, port, pass string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pass,
		DB:       0,
	})
}

func Publisher(e echo.Context) error {
	var (
		req     = new(RequestPublisher)
		err     error
		ctx     = e.Request().Context()
		payload []byte
	)
	if err = e.Bind(&req); err != nil {
		return e.JSON(400, map[string]interface{}{
			"reason":  "Bad request, check the body before send it",
			"success": false,
		})
	}

	if err = redisc(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS")).Ping(ctx).Err(); err != nil {
		return e.JSON(500, map[string]interface{}{
			"reason":  "Redis can't connect to redis-server",
			"success": false,
		})
	}

	if payload, err = json.Marshal(req); err != nil {
		return e.JSON(500, map[string]interface{}{
			"reason":  "Redis can't marshaling payload",
			"success": false,
		})
	}

	if err = redisc(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS")).Publish(ctx, os.Getenv("REDIS_CHANNEL_MESSAGE"), payload).Err(); err != nil {
		return e.JSON(500, map[string]interface{}{
			"reason":  "Redis can't publish message",
			"success": false,
		})
	}

	return e.JSON(200, map[string]interface{}{
		"reason":  "Redis sent",
		"success": true,
	})

}
