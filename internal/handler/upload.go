package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/SteGG200/storage/internal/fsutil"
)

const maxUploadSize int64 = 10 * 1024 * 1024 * 1024 // 10 GB

// sendSSEProgress sends a progress JSON via SSE.
func sendSSEProgress(w http.ResponseWriter, flusher http.Flusher, read, total int64, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	data, _ := json.Marshal(map[string]any{
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
	// #nosec G706
	log.Printf("[ERROR] 500 SSE Upload Error: %q", msg)
	mu.Lock()
	defer mu.Unlock()

	data, _ := json.Marshal(ErrorResponse{Error: msg})
	_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
	if flusher != nil {
		flusher.Flush()
	}
}

// PrepareUpload handles POST /upload/{path...} to initialize an upload session.
func (h *Handler) PrepareUpload(w http.ResponseWriter, r *http.Request) {
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

	// Limit request body to 10MB to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil && err != http.ErrNotMultipart {
		sendError(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	name := r.FormValue("name")
	if name == "" {
		sendError(w, http.StatusBadRequest, "name field is required")
		return
	}

	sizeStr := r.FormValue("size")
	if sizeStr == "" {
		sendError(w, http.StatusBadRequest, "size field is required")
		return
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil || size <= 0 {
		sendError(w, http.StatusBadRequest, "invalid or non-positive file size")
		return
	}

	if size > maxUploadSize {
		sendError(w, http.StatusRequestEntityTooLarge, "file size exceeds 10GB limit")
		return
	}

	var modTime *time.Time
	modTimeStr := r.FormValue("modifiedAt")
	if modTimeStr != "" {
		parsedTime, err := time.Parse(time.RFC3339, modTimeStr)
		if err != nil {
			parsedTime, err = time.Parse(time.RFC3339Nano, modTimeStr)
		}
		if err != nil {
			sendError(w, http.StatusBadRequest, "invalid ISO 8601 modifiedAt timestamp format")
			return
		}
		modTime = &parsedTime
	}

	if err := fsutil.ValidateName(name); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := fsutil.CheckDuplicate(targetPath, name); err != nil {
		sendError(w, http.StatusConflict, err.Error())
		return
	}

	unlock := h.Locks.Lock(targetPath)
	tempFile, err := os.CreateTemp(targetPath, "upload-*.tmp")
	if err != nil {
		unlock()
		sendError(w, http.StatusInternalServerError, "error creating temporary file")
		return
	}
	_ = tempFile.Close()
	unlock()

	filePath := filepath.Join(targetPath, name)
	h.uploadStore.Register(filePath, tempFile.Name(), size, modTime)

	sendJSON(w, http.StatusOK, map[string]string{"status": "oke"})
}

// StreamUpload handles PATCH /upload/{path...} to stream raw file content.
func (h *Handler) StreamUpload(w http.ResponseWriter, r *http.Request) {
	relPath := r.PathValue("path")

	filePath, err := fsutil.ValidatePath(h.StorageRoot, relPath)
	if err != nil {
		sendError(w, http.StatusForbidden, err.Error())
		return
	}

	tracker, ok := h.uploadStore.Tracker(filePath)
	if !ok {
		sendError(w, http.StatusNotFound, "upload target not found or expired")
		return
	}
	defer tracker.Cleanup()

	contentLengthStr := r.Header.Get("Content-Length")
	if contentLengthStr == "" {
		sendError(w, http.StatusBadRequest, "Content-Length header is required")
		return
	}

	contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil || contentLength <= 0 {
		sendError(w, http.StatusBadRequest, "invalid Content-Length")
		return
	}

	if contentLength > maxUploadSize {
		sendError(w, http.StatusRequestEntityTooLarge, "file size exceeds 10GB limit")
		return
	}

	if contentLength != tracker.totalSize {
		sendError(w, http.StatusBadRequest, "Content-Length does not match declared upload size")
		return
	}

	// #nosec G304
	tempFile, err := os.OpenFile(tracker.tempPath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "error opening temporary file")
		return
	}

	pw := &progressWriter{
		totalWritten: 0,
		targetWriter: tempFile,
		onProgress: func(written int64) {
			tracker.Publish(written)
		},
	}

	limitedBody := io.LimitReader(r.Body, tracker.totalSize)
	_, err = io.Copy(pw, limitedBody)
	if err != nil {
		_ = tempFile.Close()
		sendError(w, http.StatusInternalServerError, "error saving file")
		return
	}

	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		sendError(w, http.StatusInternalServerError, "error syncing file")
		return
	}

	if err := tempFile.Close(); err != nil {
		sendError(w, http.StatusInternalServerError, "error closing temporary file")
		return
	}

	parentDir := filepath.Dir(filePath)
	unlock := h.Locks.Lock(parentDir)
	//#nosec G304 G703
	err = os.Rename(tracker.tempPath, filePath)
	unlock()

	if err != nil {
		sendError(w, http.StatusInternalServerError, "error moving file to target location")
		return
	}

	if tracker.modTime != nil {
		//#nosec G703
		_ = os.Chtimes(filePath, *tracker.modTime, *tracker.modTime)
	}

	sendJSON(w, http.StatusOK, map[string]string{"status": "oke"})
}

// UploadProgress handles GET /progress/{path...} with SSE progress reporting.
func (h *Handler) UploadProgress(w http.ResponseWriter, r *http.Request) {
	relPath := r.PathValue("path")

	targetPath, err := fsutil.ValidatePath(h.StorageRoot, relPath)
	if err != nil {
		sendError(w, http.StatusForbidden, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		sendError(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	progressCh, unsubscribe, totalSize, ok := h.uploadStore.Subscribe(targetPath)
	if !ok {
		sendSSEError(w, flusher, "file uploading progress not found", &sync.Mutex{})
		return
	}
	defer unsubscribe()

	var mu sync.Mutex
	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return
		case byteWritten, ok := <-progressCh:
			if !ok {
				return
			}
			sendSSEProgress(w, flusher, byteWritten, totalSize, &mu)
		}
	}
}
