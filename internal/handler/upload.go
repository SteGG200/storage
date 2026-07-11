package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SteGG200/storage/internal/fsutil"
)

// sendSSEProgress sends a progress JSON via SSE.
func sendSSEProgress(w http.ResponseWriter, flusher http.Flusher, read, total int64, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	percent := int(0)
	if total > 0 {
		percent = int((read * 100) / total)
		if percent > 100 {
			percent = 100
		}
	}

	data, _ := json.Marshal(map[string]interface{}{
		"percent":      percent,
		"bytesWritten": read,
		"totalBytes":   total,
	})
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	if flusher != nil {
		flusher.Flush()
	}
}

// sendSSEError sends an error message via SSE.
func sendSSEError(w http.ResponseWriter, flusher http.Flusher, msg string, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	data, _ := json.Marshal(map[string]string{"error": msg})
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	if flusher != nil {
		flusher.Flush()
	}
}

// UploadFile handles POST /upload/{path...} with SSE progress reporting.
func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	relPath := r.PathValue("path")

	targetPath, err := fsutil.ValidatePath(h.StorageRoot, relPath)
	if err != nil {
		sendError(w, http.StatusForbidden, err.Error())
		return
	}

	if err := fsutil.RejectSymlinksAndHardlinks(targetPath); err != nil {
		sendError(w, http.StatusForbidden, err.Error())
		return
	}

	// #nosec
	fi, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(w, http.StatusNotFound, "path not found")
			return
		}
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !fi.IsDir() {
		sendError(w, http.StatusBadRequest, "target path must be a directory")
		return
	}

	// Set headers for SSE stream
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	flusher, _ := w.(http.Flusher)
	if flusher != nil {
		flusher.Flush()
	}

	var writeMu sync.Mutex
	totalBytes := r.ContentLength

	// Wrap r.Body to track read progress
	progressR := &progressReader{
		r:     r.Body,
		total: totalBytes,
		onProgress: func(read int64, total int64) {
			sendSSEProgress(w, flusher, read, total, &writeMu)
		},
	}
	r.Body = &progressReadCloser{Reader: progressR, Closer: r.Body}

	mr, err := r.MultipartReader()
	if err != nil {
		sendSSEError(w, flusher, "failed to get multipart reader: "+err.Error(), &writeMu)
		return
	}

	var fileName string
	var tempFile *os.File
	var tempPath string

	defer func() {
		if tempFile != nil {
			_ = tempFile.Close()
		}
		if tempPath != "" {
			// #nosec
			_ = os.Remove(tempPath)
		}
	}()

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			sendSSEError(w, flusher, "failed to read part: "+err.Error(), &writeMu)
			return
		}

		formName := part.FormName()
		switch formName {
		case "name":
			buf := new(strings.Builder)
			_, _ = io.Copy(buf, part)
			fileName = strings.TrimSpace(buf.String())
		case "file":
			// We can check if filename was passed in the part header if name field is empty,
			// but we prioritize form field "name".
			if fileName == "" {
				fileName = strings.TrimSpace(part.FileName())
			}

			// Create a temporary file to stream file content safely
			tempFile, err = os.CreateTemp(targetPath, ".upload-tmp-*")
			if err != nil {
				sendSSEError(w, flusher, "failed to create temporary file: "+err.Error(), &writeMu)
				return
			}
			tempPath = tempFile.Name()

			_, err = io.Copy(tempFile, part)
			if err != nil {
				sendSSEError(w, flusher, "failed to write file: "+err.Error(), &writeMu)
				return
			}
			if err := tempFile.Close(); err != nil {
				sendSSEError(w, flusher, "failed to close temporary file: "+err.Error(), &writeMu)
				return
			}
			tempFile = nil
		}
	}

	if err := fsutil.ValidateName(fileName); err != nil {
		sendSSEError(w, flusher, err.Error(), &writeMu)
		return
	}

	// Lock the parent directory during critical write section
	unlock := h.Locks.Lock(targetPath)
	defer unlock()

	if err := fsutil.CheckDuplicate(targetPath, fileName); err != nil {
		sendSSEError(w, flusher, err.Error(), &writeMu)
		return
	}

	finalPath := filepath.Join(targetPath, fileName)
	// #nosec
	if err := os.Rename(tempPath, finalPath); err != nil {
		sendSSEError(w, flusher, "failed to move uploaded file: "+err.Error(), &writeMu)
		return
	}
	tempPath = "" // Prevent defer from deleting it

	// Report final completion progress
	sendSSEProgress(w, flusher, totalBytes, totalBytes, &writeMu)

	// Send final status oke
	successData, _ := json.Marshal(map[string]interface{}{
		"status": "oke",
		"file":   fileName,
	})
	writeMu.Lock()
	_, _ = fmt.Fprintf(w, "data: %s\n\n", successData)
	if flusher != nil {
		flusher.Flush()
	}
	writeMu.Unlock()
}
