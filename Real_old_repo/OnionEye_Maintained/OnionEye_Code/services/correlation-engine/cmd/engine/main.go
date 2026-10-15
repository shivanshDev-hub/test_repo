package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/onioneye/correlation-engine/internal/handler"
)

func main() {
	port := os.Getenv("CORRELATION_PORT")
	if port == "" {
		port = "8084" // Using 8084 since 8081=ingestion, 8082=storage, 8083=content-filter
	}

	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "correlation-engine"})
	})

	// Register Engine 1 and Engine 2 stateless HTTP endpoints
	handler.RegisterRoutes(r)

	// In a real environment, this service would also subscribe to
	// `scan.completed` events on Redis/RabbitMQ and immediately re-evaluate
	// the actor's graph clustering, emitting alerts based on Composite Confidence.

	log.Printf("correlation-engine (Phase 9) listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
