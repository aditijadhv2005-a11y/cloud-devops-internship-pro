package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type response struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

func main() {
	version := env("APP_VERSION", "dev")
	failHealth := env("FAIL_HEALTH", "false") == "true"
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if failHealth {
			http.Error(w, "intentional test failure", http.StatusServiceUnavailable)
			return
		}
		writeJSON(w, http.StatusOK, response{Status: "ok", Version: version, Timestamp: time.Now().UTC().Format(time.RFC3339)})
	})
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"version": version})
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, response{Status: "hello", Version: version, Timestamp: time.Now().UTC().Format(time.RFC3339)})
	})
	srv := &http.Server{Addr: ":8080", Handler: logging(mux), ReadHeaderTimeout: 5 * time.Second}
	slog.Info("starting server", "addr", srv.Addr, "version", version)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
	})
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
