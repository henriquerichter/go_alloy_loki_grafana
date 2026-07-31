package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type response struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		message := "ok"
		now := time.Now().UTC()
		timestamp := now.Format(time.RFC3339Nano)
		writeLogAt(timestamp, now.UnixNano(), "info", message, map[string]any{"path": r.URL.Path})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{Message: message, Time: timestamp})
	})

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		message := "hello from Go"
		now := time.Now().UTC()
		timestamp := now.Format(time.RFC3339Nano)
		writeLogAt(timestamp, now.UnixNano(), "info", message, map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response{Message: message, Time: timestamp})
	})

	server := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	writeLog("info", "api server started", map[string]any{"address": server.Addr})
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		writeLog("error", "api server stopped", map[string]any{"error": err.Error()})
		os.Exit(1)
	}
}
