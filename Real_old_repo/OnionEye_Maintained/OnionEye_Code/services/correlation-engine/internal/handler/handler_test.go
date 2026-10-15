package handler_test

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/gin-gonic/gin"
"github.com/onioneye/correlation-engine/internal/engine"
"github.com/onioneye/correlation-engine/internal/handler"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler.RegisterRoutes(r)
	return r
}

func TestHandleEvaluateStylometry(t *testing.T) {
	r := setupTestRouter()

	reqPayload := handler.EvaluateStylometryRequest{
		VectorA:      []float64{0.5, 0.5, 0.5, 0.5},
		VectorB:      []float64{0.5, 0.5, 0.5, 0.5},
		SampleWordsA: 200,
		SampleWordsB: 220,
	}
	body, _ := json.Marshal(reqPayload)

	req := httptest.NewRequest(http.MethodPost, "/evaluate/stylometry", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", w.Code, w.Body.String())
	}

	var res engine.StylometryMatchResult
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.SampleGate {
		t.Fatalf("expected sample gate to pass")
	}
	if res.Similarity < 0.99 {
		t.Fatalf("expected similarity ~1.0, got %f", res.Similarity)
	}
}

func TestHandleEvaluateVector(t *testing.T) {
	r := setupTestRouter()

	art := engine.ExtractedArtifact{
		Type:       engine.VectorEmail,
		Value:      "admin@alphabay.onion",
		ObservedAt: time.Now(),
	}
	body, _ := json.Marshal(art)

	req := httptest.NewRequest(http.MethodPost, "/evaluate/vector", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var res engine.ScoredArtifact
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.HardScore > 0.20 {
		t.Fatalf("expected low hard score for generic admin email, got %f", res.HardScore)
	}
}
