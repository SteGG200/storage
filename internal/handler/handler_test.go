package handler_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/SteGG200/storage/internal/handler"
	"github.com/SteGG200/storage/internal/model"
)

// multipartField creates a multipart form body with a single field.
func multipartField(t *testing.T, fieldName, fieldValue string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField(fieldName, fieldValue); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return body, writer.FormDataContentType()
}

func setupTestHandler(t *testing.T, root string) http.Handler {
	h, err := handler.NewHandler(root)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /src/{path...}", h.GetSrc)
	mux.HandleFunc("POST /src/{path...}", h.PostSrc)
	mux.HandleFunc("PUT /src/{path...}", h.PutSrc)
	mux.HandleFunc("DELETE /src/{path...}", h.DeleteSrc)
	mux.HandleFunc("POST /upload/{path...}", h.UploadFile)
	mux.HandleFunc("GET /download/{path...}", h.DownloadFile)

	return mux
}

func TestGetSrcList(t *testing.T) {
	tmp := t.TempDir()

	// Setup hierarchy
	// file1.txt
	// dir1/
	//   file2.txt

	if err := os.WriteFile(filepath.Join(tmp, "file1.txt"), []byte("content1"), 0600); err != nil {
		t.Fatal(err)
	}
	dir1 := filepath.Join(tmp, "dir1")
	if err := os.Mkdir(dir1, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir1, "file2.txt"), []byte("content2"), 0600); err != nil {
		t.Fatal(err)
	}

	app := setupTestHandler(t, tmp)

	// Test GET /src/
	req := httptest.NewRequest("GET", "/src/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rr.Code)
	}

	var items []model.Item
	if err := json.Unmarshal(rr.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items, got: %d", len(items))
	}

	// Test GET /src/dir1
	req = httptest.NewRequest("GET", "/src/dir1", nil)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rr.Code)
	}

	var subItems []model.Item
	if err := json.Unmarshal(rr.Body.Bytes(), &subItems); err != nil {
		t.Fatal(err)
	}

	if len(subItems) != 1 || subItems[0].Name != "file2.txt" {
		t.Fatalf("expected 1 item file2.txt, got: %+v", subItems)
	}
}

func TestGetSrcSearch(t *testing.T) {
	tmp := t.TempDir()

	// Setup hierarchy
	// match_1.txt
	// other.txt
	// nested/
	//   match_2.txt

	if err := os.WriteFile(filepath.Join(tmp, "match_1.txt"), []byte("match"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "other.txt"), []byte("other"), 0600); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(tmp, "nested")
	if err := os.Mkdir(nested, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "match_2.txt"), []byte("nested match"), 0600); err != nil {
		t.Fatal(err)
	}

	app := setupTestHandler(t, tmp)

	req := httptest.NewRequest("GET", "/src/?q=match", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", rr.Code)
	}

	var items []model.Item
	if err := json.Unmarshal(rr.Body.Bytes(), &items); err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 matches, got: %d", len(items))
	}
}

func TestPostSrcCreateDir(t *testing.T) {
	tmp := t.TempDir()
	app := setupTestHandler(t, tmp)

	// Create directory under root
	body, contentType := multipartField(t, "newName", "newfolder")
	req := httptest.NewRequest("POST", "/src/", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got: %d, body: %s", rr.Code, rr.Body.String())
	}

	var res map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}

	if res["path"] != "newfolder" {
		t.Fatalf("expected path 'newfolder', got: %s", res["path"])
	}

	// Verify on disk
	fi, err := os.Stat(filepath.Join(tmp, "newfolder"))
	if err != nil || !fi.IsDir() {
		t.Fatalf("directory was not created or is not a directory: %v", err)
	}

	// Try creating duplicate folder
	body2, contentType2 := multipartField(t, "newName", "newfolder")
	req = httptest.NewRequest("POST", "/src/", body2)
	req.Header.Set("Content-Type", contentType2)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got: %d", rr.Code)
	}
}

func TestPutSrcRename(t *testing.T) {
	tmp := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmp, "oldfile.txt"), []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(tmp, "olddir"), 0700); err != nil {
		t.Fatal(err)
	}

	app := setupTestHandler(t, tmp)

	// Rename file
	body, contentType := multipartField(t, "newName", "newfile.txt")
	req := httptest.NewRequest("PUT", "/src/oldfile.txt", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, body: %s", rr.Code, rr.Body.String())
	}

	var fileRes map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &fileRes); err != nil {
		t.Fatal(err)
	}
	if fileRes["status"] != "oke" {
		t.Fatalf("expected status 'oke', got: %s", fileRes["status"])
	}

	// Verify on disk
	if _, err := os.Stat(filepath.Join(tmp, "newfile.txt")); err != nil {
		t.Fatal("renamed file not found")
	}

	// Rename folder
	body2, contentType2 := multipartField(t, "newName", "newdir")
	req = httptest.NewRequest("PUT", "/src/olddir", body2)
	req.Header.Set("Content-Type", contentType2)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got: %d", rr.Code)
	}

	var dirRes map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &dirRes); err != nil {
		t.Fatal(err)
	}
	if dirRes["path"] != "newdir" {
		t.Fatalf("expected directory rename path 'newdir', got: %s", dirRes["path"])
	}
}

func TestDeleteSrc(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "delete_me.txt")
	if err := os.WriteFile(file, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}

	app := setupTestHandler(t, tmp)

	req := httptest.NewRequest("DELETE", "/src/delete_me.txt", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got: %d", rr.Code)
	}

	if _, err := os.Stat(file); !os.IsNotExist(err) {
		t.Fatalf("file was not deleted")
	}
}

func TestDownloadFile(t *testing.T) {
	tmp := t.TempDir()
	file := filepath.Join(tmp, "download.txt")
	content := "file content here"
	if err := os.WriteFile(file, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	app := setupTestHandler(t, tmp)

	req := httptest.NewRequest("GET", "/download/download.txt", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got: %d", rr.Code)
	}

	if rr.Body.String() != content {
		t.Fatalf("expected body %q, got %q", content, rr.Body.String())
	}

	// Try downloading directory
	req = httptest.NewRequest("GET", "/download/", nil)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when downloading directory, got: %d", rr.Code)
	}
}

func TestUploadFile(t *testing.T) {
	tmp := t.TempDir()
	app := setupTestHandler(t, tmp)

	// Prepare multipart upload body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Field 'name'
	err := writer.WriteField("name", "uploaded.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Field 'file'
	part, err := writer.CreateFormFile("file", "uploaded.txt")
	if err != nil {
		t.Fatal(err)
	}
	fileContent := "hello multipart upload"
	_, _ = part.Write([]byte(fileContent))
	_ = writer.Close()

	req := httptest.NewRequest("POST", "/upload/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, body: %s", rr.Code, rr.Body.String())
	}

	// Parse SSE response
	responseLines := strings.Split(rr.Body.String(), "\n")
	var hasProgress bool
	var hasSuccess bool

	for _, line := range responseLines {
		if strings.HasPrefix(line, "data: ") {
			dataStr := strings.TrimPrefix(line, "data: ")
			var progress map[string]interface{}
			if err := json.Unmarshal([]byte(dataStr), &progress); err != nil {
				continue
			}

			// Check progress percentages
			if _, ok := progress["percent"]; ok {
				hasProgress = true
			}
			if status, ok := progress["status"]; ok && status == "oke" {
				hasSuccess = true
			}
		}
	}

	if !hasProgress {
		t.Error("expected SSE response to contain progress reports")
	}
	if !hasSuccess {
		t.Error("expected SSE response to contain success status 'oke'")
	}

	// Verify disk file
	diskPath := filepath.Join(tmp, "uploaded.txt")
	// #nosec G304
	diskData, err := os.ReadFile(diskPath)
	if err != nil {
		t.Fatalf("file not found on disk: %v", err)
	}

	if string(diskData) != fileContent {
		t.Fatalf("expected file content %q, got %q", fileContent, string(diskData))
	}
}
