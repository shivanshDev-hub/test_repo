package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// maxStoredNotifications caps the in-memory ring so a burst of worker rotation
// events can never grow the dashboard process unboundedly.
const maxStoredNotifications = 50

// notificationStore keeps the most recent notifications (worker rotations, ...)
// posted by other services, newest first. It is purely in-memory: surviving a
// dashboard restart is NOT a requirement — the popup is ephemeral and the
// crawler-worker keeps its own rotation log on stdout.
type notificationStore struct {
	mu    sync.Mutex
	items []map[string]any
}

func newNotificationStore() *notificationStore {
	return &notificationStore{items: make([]map[string]any, 0, maxStoredNotifications)}
}

// push stores a notification, guaranteeing an `at` timestamp in RFC3339Nano.
func (s *notificationStore) push(n map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if n == nil {
		n = map[string]any{}
	}
	if _, ok := n["at"]; !ok {
		n["at"] = time.Now().UTC().Format(time.RFC3339Nano)
	}

	s.items = append([]map[string]any{n}, s.items...)
	if len(s.items) > maxStoredNotifications {
		s.items = s.items[:maxStoredNotifications]
	}
}

// recent returns a copy of the stored notifications, newest first.
func (s *notificationStore) recent() []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]map[string]any, len(s.items))
	copy(out, s.items)
	return out
}

// notifyHandler accepts a best-effort notification POSTed by another service
// (e.g. crawler-worker rotation: {"event":"worker-rotation",...}), stores it for
// the frontend to poll, and echoes the payload as an ack. It never fails the
// caller: an empty/invalid body is still stored (with just an `at` stamp).
func notifyHandler(store *notificationStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload map[string]any
		if err := c.ShouldBindJSON(&payload); err != nil {
			payload = map[string]any{}
		}
		store.push(payload)
		c.JSON(http.StatusOK, gin.H{"status": "ok", "received": payload})
	}
}

// notificationsHandler returns the recent notifications, newest first, so the
// frontend can poll for push-style events like worker rotations without needing
// a websocket/SSE channel.
func notificationsHandler(store *notificationStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"items": store.recent()})
	}
}
