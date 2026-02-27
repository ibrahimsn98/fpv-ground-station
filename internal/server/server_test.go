package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"fpv-ground-station/internal/ltm"
	"fpv-ground-station/internal/telemetry"

	"nhooyr.io/websocket"
)

func testServer(t *testing.T) (*Server, *telemetry.Store, *telemetry.Stats) {
	t.Helper()

	store := &telemetry.Store{}
	stats := telemetry.NewStats()

	webFS := fstest.MapFS{
		"index.html":     &fstest.MapFile{Data: []byte("<html>FPV Ground Station</html>")},
		"assets/app.js":  &fstest.MapFile{Data: []byte("console.log('ok')")},
	}

	srv := New(Config{
		Store: store,
		Stats: stats,
		WebFS: webFS,
	})

	return srv, store, stats
}

func TestSPAHandler_ServesIndexHTML(t *testing.T) {
	srv, _, _ := testServer(t)
	handler := srv.spaHandler()

	// Request root
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET / status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "FPV Ground Station") {
		t.Fatalf("GET / body = %q, want 'FPV Ground Station'", body)
	}
}

func TestSPAHandler_ServesStaticFile(t *testing.T) {
	srv, _, _ := testServer(t)
	handler := srv.spaHandler()

	req := httptest.NewRequest("GET", "/assets/app.js", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /assets/app.js status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "console.log") {
		t.Fatalf("GET /assets/app.js body = %q", body)
	}
}

func TestSPAHandler_FallbackToIndex(t *testing.T) {
	srv, _, _ := testServer(t)
	handler := srv.spaHandler()

	// Request a non-existent path (SPA route)
	req := httptest.NewRequest("GET", "/dashboard/settings", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /dashboard/settings status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "FPV Ground Station") {
		t.Fatalf("SPA fallback should serve index.html, got %q", body)
	}
}

func TestWebSocket_ConnectAndReceive(t *testing.T) {
	srv, store, _ := testServer(t)

	// Seed some telemetry data
	store.Update(ltm.Frame{
		Function: ltm.FuncAttitude,
		Time:     time.Now(),
		Attitude: &ltm.AttitudeData{
			Roll:    12,
			Pitch:   -3,
			Heading: 270,
		},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", srv.handleWebSocket)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.broadcastLoop(ctx)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"
	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("ws dial: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	// Read one message
	readCtx, readCancel := context.WithTimeout(ctx, 2*time.Second)
	defer readCancel()

	_, data, err := conn.Read(readCtx)
	if err != nil {
		t.Fatalf("ws read: %v", err)
	}

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}

	if msg.Timestamp == 0 {
		t.Error("timestamp should be non-zero")
	}
	if msg.Attitude == nil {
		t.Fatal("attitude should be present")
	}
	if msg.Attitude.Roll != 12 {
		t.Errorf("roll = %v, want 12", msg.Attitude.Roll)
	}
	if msg.Stats == nil {
		t.Fatal("stats should be present")
	}
}

func TestWebSocket_MultipleClients(t *testing.T) {
	srv, store, _ := testServer(t)

	store.Update(ltm.Frame{
		Function: ltm.FuncStatus,
		Time:     time.Now(),
		Status:   &ltm.StatusData{Vbat: 11.8, Armed: true, FlightMode: 2},
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", srv.handleWebSocket)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go srv.broadcastLoop(ctx)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	// Connect two clients
	conn1, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("ws dial 1: %v", err)
	}
	defer conn1.Close(websocket.StatusNormalClosure, "")

	conn2, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("ws dial 2: %v", err)
	}
	defer conn2.Close(websocket.StatusNormalClosure, "")

	// Both should receive a message
	readCtx, readCancel := context.WithTimeout(ctx, 2*time.Second)
	defer readCancel()

	_, data1, err := conn1.Read(readCtx)
	if err != nil {
		t.Fatalf("ws read 1: %v", err)
	}

	_, data2, err := conn2.Read(readCtx)
	if err != nil {
		t.Fatalf("ws read 2: %v", err)
	}

	var msg1, msg2 Message
	json.Unmarshal(data1, &msg1)
	json.Unmarshal(data2, &msg2)

	if msg1.Status == nil || msg2.Status == nil {
		t.Fatal("both clients should receive status data")
	}
}
