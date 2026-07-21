package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/SteGG200/storage/internal/fsutil"
)

// sendSSEProgress sends a progress JSON via SSE.
func sendSSEProgress(w http.ResponseWriter, flusher http.Flusher, read, total int64, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	data, _ := json.Marshal(map[string]any{
		"bytesWritten": read,
		"totalBytes":   total,
	})
	_, _ = fmt.Fprintf(w, "event:info\ndata: %s\n\n", data)
	if flusher != nil {
		flusher.Flush()
	}
}

// sendSSEError sends an error message via SSE.
func sendSSEError(w http.ResponseWriter, flusher http.Flusher, msg string, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()

	data, _ := json.Marshal(ErrorResponse{Error: msg})
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

	// #nosec G120
	fileName := r.FormValue("name")
	srcFile, srcFileHeader, err := r.FormFile("file")
	if err != nil {
		sendSSEError(w, flusher, "failed to get file from form: "+err.Error(), &writeMu)
		return
	}

	if fileName == "" {
		fileName = srcFileHeader.Filename
	}

	if err := fsutil.ValidateName(fileName); err != nil {
		sendSSEError(w, flusher, err.Error(), &writeMu)
		return
	}

	unlock := h.Locks.Lock(targetPath)
	defer unlock()

	if err := fsutil.CheckDuplicate(targetPath, fileName); err != nil {
		sendSSEError(w, flusher, err.Error(), &writeMu)
		return
	}

	// #nosec G304
	finalPath := filepath.Join(targetPath, fileName)
	//#nosec G304 G703
	dstFile, err := os.Create(finalPath)
	if err != nil {
		sendSSEError(w, flusher, "failed to create file: "+err.Error(), &writeMu)
		return
	}
	defer func() {
		_ = dstFile.Close()
	}()

	pw := &progressWriter{
		totalWritten: 0,
		totalSize:    srcFileHeader.Size,
		targetWriter: dstFile,
		onProgress: func(written int64, total int64) {
			sendSSEProgress(w, flusher, written, total, &writeMu)
		},
	}

	_, err = io.Copy(pw, srcFile)
	if err != nil {
		sendSSEError(w, flusher, "failed to write file: "+err.Error(), &writeMu)
		return
	}

	successData, _ := json.Marshal(map[string]any{
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
