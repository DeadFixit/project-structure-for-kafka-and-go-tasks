package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task3-go-k8s-analytics/internal/analytics"
	"task3-go-k8s-analytics/internal/metrics"
	"task3-go-k8s-analytics/internal/model"
	"task3-go-k8s-analytics/internal/storage"

	"github.com/gorilla/mux"
)

type Server struct {
	Engine *analytics.Engine
	Cache  *storage.Cache
	Queue  chan model.Metric
}

func New() *Server {
	window := 50
	threshold := 2.0
	engine := analytics.NewEngine(window, threshold)
	cache := storage.NewCache(getEnv("REDIS_ADDR", "localhost:6379"), getEnv("REDIS_PASSWORD", ""), 0)
	return &Server{
		Engine: engine,
		Cache:  cache,
		Queue:  make(chan model.Metric, 2048),
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := mux.NewRouter()
	router.Use(metrics.Middleware)

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	router.Handle("/metrics", metrics.Handler())

	router.HandleFunc("/ingest", s.handleIngest).Methods(http.MethodPost)
	router.HandleFunc("/analytics", s.handleAnalytics).Methods(http.MethodGet)

	srv := &http.Server{Addr: ":8081", Handler: router}

	go s.backgroundWorker(ctx)

	go func() {
		<-ctx.Done()
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctxTimeout)
		_ = s.Cache.Close()
	}()

	log.Printf("analytics server listening on %s", srv.Addr)
	return srv.ListenAndServe()
}

func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) {
	var m model.Metric
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if m.Timestamp.IsZero() {
		m.Timestamp = time.Now().UTC()
	}
	select {
	case s.Queue <- m:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("queued"))
	default:
		http.Error(w, "queue full", http.StatusServiceUnavailable)
	}
}

func (s *Server) handleAnalytics(w http.ResponseWriter, r *http.Request) {
	snapshot := s.Engine.Snapshot()
	if snapshot.Anomaly {
		metrics.AnomaliesTotal.Inc()
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

func (s *Server) backgroundWorker(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-s.Queue:
			snapshot := s.Engine.Add(m)
			s.storeInCache(ctx, snapshot)
		case <-ticker.C:
			s.storeInCache(ctx, s.Engine.Snapshot())
		}
	}
}

func (s *Server) storeInCache(ctx context.Context, snapshot analytics.Snapshot) {
	payload, _ := json.Marshal(snapshot)
	ttl := 5 * time.Minute
	_ = s.Cache.SetMetric(ctx, "analytics:last", string(payload), ttl)
}

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := New()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	return s.Start(ctx)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
