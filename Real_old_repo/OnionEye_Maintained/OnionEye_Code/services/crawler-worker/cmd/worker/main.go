package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/onioneye/crawler-worker/internal/config"
	"github.com/onioneye/crawler-worker/internal/crawler"
	"github.com/onioneye/crawler-worker/internal/worker"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}

	// Ensure the stream + consumer group exist (idempotent).
	if err := rdb.XGroupCreateMkStream(ctx, cfg.Stream, cfg.Group, "$").Err(); err != nil &&
		!isGroupExists(err) {
		log.Printf("create group: %v", err)
	}

	// Tor crawler client
	cr, err := crawler.NewClient(cfg)
	if err != nil {
		log.Fatalf("crawler client: %v", err)
	}

	log.Printf("crawler-worker starting: workers=%d stream=%s group=%s consumer=%s socks=%s",
		cfg.WorkerCount, cfg.Stream, cfg.Group, cfg.Consumer, cfg.TorSocksAddr)

	// Read-only status surface (Tor health + liveness) consumed by dashboard-api
	// for the UI HUD. Must start before pool.Run: Run blocks forever.
	pool := worker.NewPool(cfg, rdb, cr)
	startStatusServer(cfg, pool)

	pool.Run(ctx)
}

func isGroupExists(err error) bool {
	return err != nil && (redis.HasErrorPrefix(err, "BUSYGROUP") ||
		redis.HasErrorPrefix(err, "ERR.*already exists"))
}
