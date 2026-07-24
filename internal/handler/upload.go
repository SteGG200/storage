package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	// #nosec G706
	log.Printf("[ERROR] 500 SSE Upload Error: %s", msg)
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

	// 1. Parse request form and file FIRST before writing/flushing response headers
	// #nosec G120
	fileName := r.FormValue("name")
	srcFile, srcFileHeader, err := r.FormFile("file")
	if err != nil {
		sendError(w, http.StatusBadRequest, "failed to get file from form: "+err.Error())
		return
	}
	defer func() {
		_ = srcFile.Close()
	}()

	if fileName == "" {
		fileName = srcFileHeader.Filename
	}

	if err := fsutil.ValidateName(fileName); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	unlock := h.Locks.Lock(targetPath)
	defer unlock()

	if err := fsutil.CheckDuplicate(targetPath, fileName); err != nil {
		sendError(w, http.StatusConflict, err.Error())
		return
	}

	// #nosec G304
	finalPath := filepath.Join(targetPath, fileName)
	//#nosec G304 G703
	dstFile, err := os.Create(finalPath)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to create file: "+err.Error())
		return
	}
	defer func() {
		_ = dstFile.Close()
	}()

	// 2. Set headers for SSE stream AFTER parsing request form body
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
