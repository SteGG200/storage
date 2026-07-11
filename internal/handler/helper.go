package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

// ErrorResponse represents an API error structure.
type ErrorResponse struct {
	Error string `json:"error"`
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, status int, msg string) {
	sendJSON(w, status, ErrorResponse{Error: msg})
}

// progressReader wraps an io.Reader and reports progress of reads.
type progressReader struct {
	r          io.Reader
	total      int64
	read       int64
	onProgress func(read int64, total int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	if n > 0 {
		pr.read += int64(n)
		if pr.onProgress != nil {
			pr.onProgress(pr.read, pr.total)
		}
	}
	return n, err
}

// progressReadCloser wraps progressReader and implements io.ReadCloser.
type progressReadCloser struct {
	io.Reader
	io.Closer
}
