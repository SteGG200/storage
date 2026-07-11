package handler

import (
	"net/http"
	"os"

	"github.com/SteGG200/storage/internal/fsutil"
)

// DownloadFile handles GET /download/{path...}.
func (h *Handler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	relPath := r.PathValue("path")
	if relPath == "" {
		sendError(w, http.StatusBadRequest, "path cannot be empty")
		return
	}

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
			sendError(w, http.StatusNotFound, "file not found")
			return
		}
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if fi.IsDir() {
		sendError(w, http.StatusBadRequest, "cannot download a directory")
		return
	}

	// Serve target file
	http.ServeFile(w, r, targetPath)
}
