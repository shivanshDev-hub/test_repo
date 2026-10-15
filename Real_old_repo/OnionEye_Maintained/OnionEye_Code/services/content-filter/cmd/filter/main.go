package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/onioneye/content-filter/internal/classifier"
)

func main() {
	// Let Gin run in release mode if we aren't debugging.
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 1. Initialize inference backends.
	// All three vectors use the same ONNXBackend runtime. Each model file is
	// optional: if it is present it loads for real; if it is missing (or the
	// runtime cannot load it) that vector falls back to a deterministic mock
	// and logs exactly what happened. USE_MOCK_MODELS=false does NOT hard-fail
	// unavailable models — it just enables the real-ONNX attempt.
	useMocks := strings.ToLower(os.Getenv("USE_MOCK_MODELS")) != "false"

	var age, nudity, behavior classifier.Inferencer
	if useMocks {
		log.Println("Mock mode (USE_MOCK_MODELS != false): all three vectors use deterministic stubs")
		// By default return benign scores so content passes.
		age = classifier.NewMock("insightface_V2+mock", 0.99)
		nudity = classifier.NewMock("nudenet_V2+mock", 0.01)
		behavior = classifier.NewMock("pose_mlp_V1+mock", 0.05)
	} else {
		log.Println("Production mode: attempting to load real ONNX backends (fail-soft per vector)")

		agePath := getEnvOrDefault("MODEL_PATH_AGE", "/models/age_mlp.onnx")
		nudityPath := getEnvOrDefault("MODEL_PATH_NUDITY", "/models/nudity_mobilenet.onnx")
		behaviorPath := getEnvOrDefault("MODEL_PATH_BEHAVIOR", "/models/pose_mlp.onnx")

		// ── Vector A — Age MLP over a 512-d face embedding ──────────────────
		ageSpec := classifier.AgeModelSpec()
		if backend, ok := tryLoadONNX("Vector A/age", agePath, ageSpec); ok {
			age = backend
		} else {
			age = classifier.NewMock("insightface_V2+mock", 0.99)
		}

		// ── Vector B — Nudity CNN classifier ────────────────────────────────
		nuditySize := getEnvIntOrDefault("NUDITY_MODEL_SIZE", 224)
		nuditySpec, err := classifier.NudityModelSpec(nuditySize)
		if err != nil {
			log.Fatalf("bad nudity model size: %v", err)
		}
		if backend, ok := tryLoadONNX("Vector B/nudity", nudityPath, nuditySpec); ok {
			nudity = backend
		} else {
			nudity = classifier.NewMock("nudenet_V2+mock", 0.01)
		}

		// ── Vector C — Behavior pose-MLP ────────────────────────────────────
		behaviorSpec, err := classifier.BehaviorModelSpec()
		if err != nil {
			log.Fatalf("bad behavior model spec: %v", err)
		}
		if backend, ok := tryLoadONNX("Vector C/behavior", behaviorPath, behaviorSpec); ok {
			behavior = backend
		} else {
			behavior = classifier.NewMock("pose_mlp_V1+mock", 0.05)
		}
	}

	if age == nil || nudity == nil || behavior == nil {
		log.Fatal("a pipeline model was not initialized")
	}

	log.Printf("Model roster → age: %s | nudity: %s | behavior: %s",
		age.Name(), nudity.Name(), behavior.Name())

	// 2. Build the pipeline.
	pipeline, err := classifier.NewPipeline(age, nudity, behavior)
	if err != nil {
		log.Fatalf("failed to construct pipeline: %v", err)
	}

	// 3. Set up the HTTP handler and router.
	h := classifier.NewHandler(pipeline)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/healthz"},
	}))

	h.RegisterRoutes(r)

	// 4. Start the server.
	port := getEnvOrDefault("PORT", "8082")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("content-filter (Phase 8A) listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %v", err)
		}
	}()

	// 5. Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited successfully")
}

func getEnvOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

// getEnvIntOrDefault parses an integer env var, falling back to def on empty
// or unparsable values (machine must not die on a typo'd knob).
func getEnvIntOrDefault(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
		log.Printf("WARNING: %s=%q is not an integer; using %d", key, val, def)
	}
	return def
}

// tryLoadONNX eagerly constructs and Load()s a real ONNX backend for modelPath
// with the given spec. It returns (backend, true) on success. On ANY failure —
// missing file, missing runtime, session creation error — it logs the reason
// and returns (nil, false) so the caller can fall back to a mock. This is the
// fail-soft contract: an unavailable model must never take the service down.
func tryLoadONNX(label, modelPath string, spec classifier.ONNXModelSpec) (classifier.Inferencer, bool) {
	if modelPath == "" {
		log.Printf("WARNING: %s: no model path configured; using mock", label)
		return nil, false
	}
	if _, err := os.Stat(modelPath); err != nil {
		log.Printf("WARNING: %s: model %s not present (%v); using mock", label, modelPath, err)
		return nil, false
	}
	backend, err := classifier.NewONNXModelBackend(label, modelPath, spec)
	if err != nil {
		log.Printf("WARNING: %s: cannot construct backend for %s: %v; using mock", label, modelPath, err)
		return nil, false
	}
	if err := backend.Load(); err != nil {
		log.Printf("WARNING: %s: ONNX session failed to load %s: %v; using mock", label, modelPath, err)
		return nil, false
	}
	log.Printf("OK: %s loaded real model from %s", label, modelPath)
	return backend, true
}
