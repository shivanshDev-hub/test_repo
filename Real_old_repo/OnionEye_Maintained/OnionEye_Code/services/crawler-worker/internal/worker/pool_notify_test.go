package worker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/onioneye/crawler-worker/internal/config"
)

// captureNotification spins up a mock dashboard API that records the request
// method + decoded JSON body, and returns either the given status code or the
// captured payload echo. It lets tests pin down the worker-rotation wire format.
func captureNotification(t *testing.T, status int) (*httptest.Server, func() *http.Request, func() map[string]any) {
	t.Helper()

	var gotReq *http.Request
	var gotBody map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReq = r
		if r.Body != nil {
			defer r.Body.Close()
			_ = json.NewDecoder(r.Body).Decode(&gotBody)
		}
		if status == http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"error":"boom"}`))
	}))

	return srv, func() *http.Request { return gotReq }, func() map[string]any { return gotBody }
}

func TestPoolNotifyRotationPayload(t *testing.T) {
	srv, req, body := captureNotification(t, http.StatusOK)
	defer srv.Close()

	p := &Pool{
		cfg: &config.Config{
			DashboardAPIURL:        srv.URL,
			DashboardNotifyPath:    "/api/v1/dashboard/notify",
			DashboardNotifyTimeout: 2 * time.Second,
		},
	}
	p.rotationID = "rot-0000000001"

	p.notifyRotation(context.Background(), "scheduled rotation")

	if req() == nil {
		t.Fatal("expected HTTP POST request to the mock dashboard, got none")
	}
	if req().Method != http.MethodPost {
		t.Fatalf("expected POST, got %s", req().Method)
	}

	b := body()
	if b == nil {
		t.Fatal("expected JSON body, got nil")
	}
	if b["event"] != "worker-rotation" {
		t.Fatalf("event = %v, want worker-rotation", b["event"])
	}
	if b["rotation_id"] != "rot-0000000001" {
		t.Fatalf("rotation_id = %v, want rot-0000000001", b["rotation_id"])
	}
	if b["reason"] != "scheduled rotation" {
		t.Fatalf("reason = %v, want scheduled rotation", b["reason"])
	}
	if _, ok := b["at"].(string); !ok {
		t.Fatalf("at missing or not a string: %T", b["at"])
	}
}

func TestPoolNotifyRotationNon2xxIsSwallowed(t *testing.T) {
	srv, req, _ := captureNotification(t, http.StatusInternalServerError)
	defer srv.Close()

	p := &Pool{
		cfg: &config.Config{
			DashboardAPIURL:        srv.URL,
			DashboardNotifyPath:    "/notify",
			DashboardNotifyTimeout: 2 * time.Second,
		},
	}
	p.rotationID = "rot-0000000002"

	// Must NOT panic / block on non-2xx — best effort only.
	p.notifyRotation(context.Background(), "scheduled rotation")

	if req() == nil {
		t.Fatal("expected the request to still be made for non-2xx responses")
	}
}

func TestPoolNotifyRotationUnreachableIsSwallowed(t *testing.T) {
	// Point at a server we close immediately: transport error path.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close()

	p := &Pool{
		cfg: &config.Config{
			DashboardAPIURL:        url,
			DashboardNotifyPath:    "/notify",
			DashboardNotifyTimeout: 500 * time.Millisecond,
		},
	}
	p.rotationID = "rot-0000000003"

	p.notifyRotation(context.Background(), "scheduled rotation") // must not panic
}

func TestPoolNotifyRotationNilConfigNoop(t *testing.T) {
	p := &Pool{} // cfg == nil
	p.notifyRotation(context.Background(), "scheduled rotation") // must not panic
}