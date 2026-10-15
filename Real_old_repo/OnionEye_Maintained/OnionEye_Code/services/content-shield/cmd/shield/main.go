package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/onioneye/content-shield/internal/handler"
	"github.com/onioneye/content-shield/internal/hashset"
	"github.com/onioneye/content-shield/internal/phash"
	"github.com/onioneye/content-shield/internal/storage"
)

func main() {
	port := env("CONTENT_SHIELD_PORT", "8082")
	hashPath := env("CSAM_HASH_SET_PATH", "/data/hashes/csam_hashes.sha256")
	storageURL := env("STORAGE_API_URL", "http://storage-api:8081")

	// Load the known-bad hash set at startup (in-memory for membership checks).
	set, err := hashset.FromFile(hashPath)
	if err != nil {
		// Demo: start with an empty set rather than crash when the volume is
		// empty — the shield still returns clean passes and logs membership 0.
		log.Printf("WARN: hash set not loaded (%v) — starting with empty set", err)
		set = hashset.New()
	}
	log.Printf("hash set loaded: %d entries", set.Len())

	h := handler.New(
		set,
		phash.NewSHA256(),
		storage.New(storageURL),
	)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "content-shield", "hash_entries": set.Len()})
	})
	r.POST("/shield", h.Handle)

	log.Printf("content-shield listening on :%s (hash_set=%s)", port, hashPath)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}