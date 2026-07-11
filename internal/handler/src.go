package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/SteGG200/storage/internal/fsutil"
	"github.com/SteGG200/storage/internal/locks"
)

// Handler handles HTTP requests for storage operations.
type Handler struct {
	StorageRoot string
	Locks       *locks.PathLocks
}

// NewHandler initializes a Handler.
func NewHandler(storageRoot string) (*Handler, error) {
	absRoot, err := filepath.Abs(storageRoot)
	if err != nil {
		return nil, err
	}
	return &Handler{
		StorageRoot: absRoot,
		Locks:       &locks.PathLocks{},
	}, nil
}

// GetSrc handles GET /src/{path...} (listing/searching).
func (h *Handler) GetSrc(w http.ResponseWriter, r *http.Request) {
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
		sendError(w, http.StatusBadRequest, "path is not a directory")
		return
	}

	q := r.URL.Query().Get("q")
	if q != "" {
		items, err := fsutil.Search(h.StorageRoot, targetPath, q)
		if err != nil {
			sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		sendJSON(w, http.StatusOK, items)
		return
	}

	items, err := fsutil.ListDir(h.StorageRoot, targetPath)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendJSON(w, http.StatusOK, items)
}

// PostSrc handles POST /src/{path...} (create folder).
func (h *Handler) PostSrc(w http.ResponseWriter, r *http.Request) {
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
		sendError(w, http.StatusBadRequest, "path must be a directory")
		return
	}

	// Limit request body to 10MB to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		sendError(w, http.StatusBadRequest, "failed to parse multipart form data")
		return
	}

	newName := r.FormValue("newName")
	if err := fsutil.ValidateName(newName); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Lock parent directory
	unlock := h.Locks.Lock(targetPath)
	defer unlock()

	if err := fsutil.CheckDuplicate(targetPath, newName); err != nil {
		sendError(w, http.StatusConflict, err.Error())
		return
	}

	newFolderPath := filepath.Join(targetPath, newName)
	// #nosec
	if err := os.Mkdir(newFolderPath, 0750); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	relNewFolder, err := filepath.Rel(h.StorageRoot, newFolderPath)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	relNewFolder = filepath.ToSlash(relNewFolder)

	sendJSON(w, http.StatusCreated, map[string]string{"path": relNewFolder})
}

// PutSrc handles PUT /src/{path...} (rename).
func (h *Handler) PutSrc(w http.ResponseWriter, r *http.Request) {
	relPath := r.PathValue("path")
	if relPath == "" {
		sendError(w, http.StatusBadRequest, "cannot rename storage root")
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
	fi, err := os.Lstat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(w, http.StatusNotFound, "path not found")
			return
		}
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Limit request body to 10MB to prevent memory exhaustion
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		sendError(w, http.StatusBadRequest, "failed to parse multipart form data")
		return
	}

	newName := r.FormValue("newName")
	if err := fsutil.ValidateName(newName); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	parentDir := filepath.Dir(targetPath)
	newPath := filepath.Join(parentDir, newName)

	// Lock parent directory
	unlock := h.Locks.Lock(parentDir)
	defer unlock()

	// #nosec
	if _, err := os.Lstat(newPath); err == nil {
		sendError(w, http.StatusConflict, "destination path already exists")
		return
	}

	// #nosec
	if err := os.Rename(targetPath, newPath); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if fi.IsDir() {
		relNew, err := filepath.Rel(h.StorageRoot, newPath)
		if err != nil {
			sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		relNew = filepath.ToSlash(relNew)
		sendJSON(w, http.StatusOK, map[string]string{"path": relNew})
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"status": "oke"})
}

// DeleteSrc handles DELETE /src/{path...} (delete).
func (h *Handler) DeleteSrc(w http.ResponseWriter, r *http.Request) {
	relPath := r.PathValue("path")
	if relPath == "" {
		sendError(w, http.StatusBadRequest, "cannot delete storage root")
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
	_, err = os.Lstat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			sendError(w, http.StatusNotFound, "path not found")
			return
		}
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	parentDir := filepath.Dir(targetPath)
	unlock := h.Locks.Lock(parentDir)
	defer unlock()

	// #nosec
	if err := os.RemoveAll(targetPath); err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{"status": "oke"})
}
