package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// ErrorResponse represents an API error structure.
type ErrorResponse struct {
	Error string `json:"error"`
}

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	if status == http.StatusInternalServerError {
		// #nosec G706
		log.Printf("[ERROR] 500 Internal Server Error: %s", msg)
	}
	sendJSON(w, status, ErrorResponse{Error: msg})
}

type progressWriter struct {
	totalWritten int64
	targetWriter io.Writer
	onProgress   func(written int64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.targetWriter.Write(p)
	if err != nil {
		return n, err
	}

	pw.totalWritten += int64(n)
	pw.onProgress(pw.totalWritten)
	return n, nil
}
