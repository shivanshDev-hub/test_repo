package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/onioneye/storage-api/internal/db"
	"github.com/onioneye/storage-api/internal/handler"
)

func main() {
	port := os.Getenv("STORAGE_API_PORT")
	if port == "" {
		port = "8081"
	}

	// This is the ONLY service that holds DATABASE_URL — every other service
	// persists through this HTTP boundary (zero-trust architecture).
	cfg := db.Load()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	h := handler.New(pool)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "storage-api"})
	})

	// Write endpoints (consumed by crawler-worker, content-shield,
	// correlation-engine).
	r.POST("/scans", h.CreateScan)
	r.POST("/content-flags", h.CreateFlag)
	r.POST("/actors", h.UpsertActor)
	r.POST("/identifiers", h.UpsertIdentifier)
	r.POST("/persona-links", h.UpsertPersonaLink)

	// Delete endpoints (consumed by dashboard-api / UI).
	r.DELETE("/delete/node", h.DeleteNode)
	r.DELETE("/delete/all", h.DeleteAll)

	// Read endpoints (consumed by dashboard-api).
	r.GET("/actors", h.ListActors)
	r.GET("/actors/:id", h.GetActor)
	r.GET("/actors/:id/behavioral", h.GetActorBehavioral)
	r.GET("/leads", h.ListLeads)
	r.GET("/stats", h.GetStats)
	r.GET("/scans", h.ListScans)
	r.GET("/images", h.ListImages)

	// Negative-control console (LeftControlRail NEGATIVE CONTROL tab).
	r.GET("/negative-control", h.ListNegativeControl)
	r.POST("/negative-control", h.CreateNegativeControlRule)
	r.DELETE("/negative-control/:id", h.DeleteNegativeControlRule)
	r.POST("/negative-control/artifacts", h.AddNegativeControlArtifact)
	r.GET("/artifacts", h.ListNegativeControlArtifacts)

	// Screenshots (render-sandbox pipeline):
	// crawler-worker POSTs captures here; dashboard-api reads metadata + PNG.
	r.POST("/screenshots", h.UpsertScreenshot)
	r.GET("/screenshots", h.GetScreenshotByURL)
	r.GET("/screenshots/:id/png", h.GetScreenshotPNG)

	log.Printf("storage-api listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
