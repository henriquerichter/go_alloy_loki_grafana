package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	alloyURL    = getenv("ALLOY_URL", "http://alloy:3101/loki/api/v1/push")
	client      = &http.Client{Timeout: 3 * time.Second}
	logger      = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	errorLogger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
)

type lokiPush struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func writeLog(level, message string, fields map[string]any) {
	now := time.Now().UTC()
	writeLogAt(now.Format(time.RFC3339Nano), now.UnixNano(), level, message, fields)
}

func writeLogAt(timestamp string, unixNano int64, level, message string, fields map[string]any) {
	entry := map[string]any{
		"time":    timestamp,
		"level":   level,
		"message": message,
	}
	for key, value := range fields {
		entry[key] = value
	}

	line, err := json.Marshal(entry)
	if err != nil {
		errorLogger.Error("encode log", "error", err)
		return
	}
	// Keep a structured local copy in the app console while also sending the same entry to Alloy.
	logger.Info("log entry", "entry", entry)

	payload := lokiPush{Streams: []lokiStream{{
		Stream: map[string]string{"app": "api", "source": "direct-http"},
		Values: [][2]string{{strconv.FormatInt(unixNano, 10), string(line)}},
	}}}
	body, err := json.Marshal(payload)
	if err != nil {
		errorLogger.Error("encode Loki payload", "error", err)
		return
	}

	request, err := http.NewRequest(http.MethodPost, alloyURL, strings.NewReader(string(body)))
	if err != nil {
		errorLogger.Error("create Alloy request", "error", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		errorLogger.Error("send log to Alloy", "error", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, response.Body)
		errorLogger.Error("Alloy returned an error status while sending log", "status", response.Status)
	}
}
