package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/onioneye/ingestion-api/internal/config"
	"github.com/onioneye/ingestion-api/internal/handler"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
	})

	h, err := handler.New(cfg, rdb)
	if err != nil {
		log.Fatalf("init handler: %v", err)
	}

	r := gin.Default()
	h.RegisterRoutes(r)

	log.Printf("ingestion-api listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
