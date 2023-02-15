package controller

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func Listener() {
	log.Println("subscriber listen")

	ctx := context.Background()

	subscriber := redisc(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASS")).Subscribe(ctx, os.Getenv("REDIS_CHANNEL_MESSAGE"))
	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(ctx)
			if err != nil {
				panic(err)
			}
			registerService(msg)
		}
	}()
}

func registerService(msg *redis.Message) {

	log.Println("msg", msg.Payload)

}
