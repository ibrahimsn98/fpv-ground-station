package server

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"fpv-ground-station/internal/telemetry"
)

// Config configures the web server.
type Config struct {
	Store    *telemetry.Store
	Stats    *telemetry.Stats
	TrackLog *telemetry.TrackLog
	Addr     string
	WebFS    fs.FS // embedded or nil in dev mode
	DevMode  bool
}

// Server serves the web UI and WebSocket telemetry.
type Server struct {
	store    *telemetry.Store
	stats    *telemetry.Stats
	trackLog *telemetry.TrackLog
	addr     string
	webFS    fs.FS
	devMode  bool

	mu      sync.RWMutex
	clients map[*client]struct{}
}

type client struct {
	send chan []byte
}

// New creates a new Server.
func New(cfg Config) *Server {
	return &Server{
		store:    cfg.Store,
		stats:    cfg.Stats,
		trackLog: cfg.TrackLog,
		addr:     cfg.Addr,
		webFS:    cfg.WebFS,
		devMode:  cfg.DevMode,
		clients:  make(map[*client]struct{}),
	}
}

// ListenAndServe starts the HTTP server and blocks until ctx is cancelled.
func (s *Server) ListenAndServe(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", s.handleWebSocket)
	mux.HandleFunc("/api/track", s.handleTrack)

	// SPA file serving (only if webFS is available)
	if s.webFS != nil {
		mux.Handle("/", s.spaHandler())
	}

	srv := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	// Start broadcast loop
	go s.broadcastLoop(ctx)

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("Web server listening on %s", s.addr)
	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// spaHandler serves static files, falling back to index.html for SPA routing.
func (s *Server) spaHandler() http.Handler {
	fileServer := http.FileServer(http.FS(s.webFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to open the file
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = path[1:] // strip leading /
		}

		f, err := s.webFS.Open(path)
		if err != nil {
			// File not found — serve index.html for SPA routing
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		f.Close()

		fileServer.ServeHTTP(w, r)
	})
}

func (s *Server) addClient(c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c] = struct{}{}
}

func (s *Server) removeClient(c *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c)
	close(c.send)
}

func (s *Server) handleTrack(w http.ResponseWriter, r *http.Request) {
	if s.trackLog == nil {
		http.Error(w, "track log not configured", http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodGet:
		points, err := s.trackLog.ReadAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(points)

	case http.MethodDelete:
		if err := s.trackLog.Clear(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
