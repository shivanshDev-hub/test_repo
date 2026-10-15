package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNotificationStorePushAndRecent(t *testing.T) {
	s := newNotificationStore()

	// Ring cap is enforced.
	for i := 0; i < maxStoredNotifications+10; i++ {
		s.push(map[string]any{"event": "worker-rotation", "n": i})
	}
	items := s.recent()
	if len(items) != maxStoredNotifications {
		t.Fatalf("expected %d stored items, got %d", maxStoredNotifications, len(items))
	}
	// Newest first.
	if items[0]["n"] != maxStoredNotifications+9 {
		t.Fatalf("expected newest item first, got %v", items[0])
	}
}

func TestNotifyEndpointStoresAndEchoes(t *testing.T) {
	store := newNotificationStore()
	r := gin.New()
	r.POST("/notify", notifyHandler(store))
	r.GET("/notifications", notificationsHandler(store))

	payload := `{"event":"worker-rotation","rotation_id":"rot-test","reason":"scheduled rotation","at":"2025-01-01T00:00:00Z"}`

	// POST the notification.
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /notify status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	var echo struct {
		Status   string         `json:"status"`
		Received map[string]any `json:"received"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &echo); err != nil {
		t.Fatalf("decode echo: %v", err)
	}
	if echo.Status != "ok" {
		t.Fatalf("echo status = %q, want ok", echo.Status)
	}
	if echo.Received["event"] != "worker-rotation" || echo.Received["rotation_id"] != "rot-test" {
		t.Fatalf("echo payload mismatch: %v", echo.Received)
	}

	// Poll it back via GET /notifications.
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("GET /notifications status = %d, want %d", w2.Code, http.StatusOK)
	}
	var out struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(w2.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode notifications: %v", err)
	}
	if len(out.Items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(out.Items))
	}
	if out.Items[0]["event"] != "worker-rotation" {
		t.Fatalf("unexpected item: %v", out.Items[0])
	}
}

func TestNotifyInvalidBodyStillAck(t *testing.T) {
	store := newNotificationStore()
	r := gin.New()
	r.POST("/notify", notifyHandler(store))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notify", bytes.NewBufferString("not-json{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("invalid body should still ack (best-effort), got status %d", w.Code)
	}
}
