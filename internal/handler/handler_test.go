package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

// multipartFields creates a multipart form body with multiple fields.
func multipartFields(t *testing.T, fields map[string]string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
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
	mux.HandleFunc("POST /upload/{path...}", h.PrepareUpload)
	mux.HandleFunc("PATCH /upload/{path...}", h.StreamUpload)
	mux.HandleFunc("GET /progress/{path...}", h.UploadProgress)
	mux.HandleFunc("GET /download/{path...}", h.DownloadFile)

	return mux
}

func TestGetSrcList(t *testing.T) {
	tmp := t.TempDir()

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

	fi, err := os.Stat(filepath.Join(tmp, "newfolder"))
	if err != nil || !fi.IsDir() {
		t.Fatalf("directory was not created or is not a directory: %v", err)
	}

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

	if _, err := os.Stat(filepath.Join(tmp, "newfile.txt")); err != nil {
		t.Fatal("renamed file not found")
	}

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

	req = httptest.NewRequest("GET", "/download/", nil)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when downloading directory, got: %d", rr.Code)
	}
}

func TestPrepareUpload(t *testing.T) {
	tmp := t.TempDir()
	app := setupTestHandler(t, tmp)

	body, contentType := multipartFields(t, map[string]string{
		"name": "uploaded.txt",
		"size": "22",
	})
	req := httptest.NewRequest("POST", "/upload/", body)
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, body: %s", rr.Code, rr.Body.String())
	}

	var res map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if res["status"] != "oke" {
		t.Fatalf("expected status 'oke', got: %s", res["status"])
	}

	body2, contentType2 := multipartFields(t, map[string]string{"size": "22"})
	req = httptest.NewRequest("POST", "/upload/", body2)
	req.Header.Set("Content-Type", contentType2)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing name, got: %d", rr.Code)
	}

	body3, contentType3 := multipartFields(t, map[string]string{"name": "test.txt"})
	req = httptest.NewRequest("POST", "/upload/", body3)
	req.Header.Set("Content-Type", contentType3)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing size, got: %d", rr.Code)
	}

	body4, contentType4 := multipartFields(t, map[string]string{"name": "test.txt", "size": "0"})
	req = httptest.NewRequest("POST", "/upload/", body4)
	req.Header.Set("Content-Type", contentType4)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid size, got: %d", rr.Code)
	}

	body5, contentType5 := multipartFields(t, map[string]string{"name": "huge.txt", "size": "10737418241"})
	req = httptest.NewRequest("POST", "/upload/", body5)
	req.Header.Set("Content-Type", contentType5)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 for size > 10GB, got: %d", rr.Code)
	}

	// Test invalid modifiedAt
	body6, contentType6 := multipartFields(t, map[string]string{
		"name":       "invalid_date.txt",
		"size":       "10",
		"modifiedAt": "not-a-date",
	})
	req = httptest.NewRequest("POST", "/upload/", body6)
	req.Header.Set("Content-Type", contentType6)
	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid modifiedAt, got: %d", rr.Code)
	}
}

func TestStreamUpload(t *testing.T) {
	tmp := t.TempDir()
	app := setupTestHandler(t, tmp)

	fileContent := "hello streaming upload"
	fileSize := fmt.Sprintf("%d", len(fileContent))

	prepBody, prepCT := multipartFields(t, map[string]string{
		"name": "streamed.txt",
		"size": fileSize,
	})
	req := httptest.NewRequest("POST", "/upload/", prepBody)
	req.Header.Set("Content-Type", prepCT)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("prepare failed, status: %d", rr.Code)
	}

	patchReq := httptest.NewRequest("PATCH", "/upload/streamed.txt", strings.NewReader(fileContent))
	patchReq.Header.Set("Content-Length", fileSize)
	patchReq.Header.Set("Content-Type", "application/octet-stream")
	patchRr := httptest.NewRecorder()
	app.ServeHTTP(patchRr, patchReq)

	if patchRr.Code != http.StatusOK {
		t.Fatalf("stream upload failed, status: %d, body: %s", patchRr.Code, patchRr.Body.String())
	}

	diskPath := filepath.Join(tmp, "streamed.txt")
	// #nosec G304
	diskData, err := os.ReadFile(diskPath)
	if err != nil {
		t.Fatalf("file not found on disk: %v", err)
	}
	if string(diskData) != fileContent {
		t.Fatalf("expected content %q, got %q", fileContent, string(diskData))
	}
}

func TestStreamUploadWithModifiedAt(t *testing.T) {
	tmp := t.TempDir()
	app := setupTestHandler(t, tmp)

	fileContent := "hello modified timestamp"
	fileSize := fmt.Sprintf("%d", len(fileContent))
	customTimeStr := "2026-07-28T10:00:00Z"
	expectedTime, err := time.Parse(time.RFC3339, customTimeStr)
	if err != nil {
		t.Fatal(err)
	}

	prepBody, prepCT := multipartFields(t, map[string]string{
		"name":       "mod_test.txt",
		"size":       fileSize,
		"modifiedAt": customTimeStr,
	})
	req := httptest.NewRequest("POST", "/upload/", prepBody)
	req.Header.Set("Content-Type", prepCT)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("prepare failed, status: %d", rr.Code)
	}

	patchReq := httptest.NewRequest("PATCH", "/upload/mod_test.txt", strings.NewReader(fileContent))
	patchReq.Header.Set("Content-Length", fileSize)
	patchReq.Header.Set("Content-Type", "application/octet-stream")
	patchRr := httptest.NewRecorder()
	app.ServeHTTP(patchRr, patchReq)

	if patchRr.Code != http.StatusOK {
		t.Fatalf("stream upload failed, status: %d, body: %s", patchRr.Code, patchRr.Body.String())
	}

	diskPath := filepath.Join(tmp, "mod_test.txt")
	fi, err := os.Stat(diskPath)
	if err != nil {
		t.Fatalf("file not found on disk: %v", err)
	}

	if !fi.ModTime().Equal(expectedTime) {
		t.Fatalf("expected modTime %v, got %v", expectedTime, fi.ModTime())
	}
}

func TestUploadProgress(t *testing.T) {
	tmp := t.TempDir()
	app := setupTestHandler(t, tmp)

	ts := httptest.NewServer(app)
	defer ts.Close()

	fileContent := "01234567890123456789012345678901234567890123456789"
	fileSize := fmt.Sprintf("%d", len(fileContent))

	prepBody, prepCT := multipartFields(t, map[string]string{
		"name": "progress_test.txt",
		"size": fileSize,
	})
	resp, err := http.Post(ts.URL+"/upload/", prepCT, prepBody)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("prepare failed: %d", resp.StatusCode)
	}

	sseClient := ts.Client()
	progressReq, err := http.NewRequest("GET", ts.URL+"/progress/progress_test.txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	sseResp, err := sseClient.Do(progressReq)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = sseResp.Body.Close()
	}()

	if sseResp.StatusCode != http.StatusOK {
		t.Fatalf("progress failed: %d", sseResp.StatusCode)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
		patchReq, _ := http.NewRequest("PATCH", ts.URL+"/upload/progress_test.txt", strings.NewReader(fileContent))
		patchReq.Header.Set("Content-Length", fileSize)
		patchReq.Header.Set("Content-Type", "application/octet-stream")
		patchResp, err := http.DefaultClient.Do(patchReq)
		if err == nil {
			_ = patchResp.Body.Close()
		}
	}()

	buf := make([]byte, 1024)
	n, _ := io.ReadAtLeast(sseResp.Body, buf, 1)
	sseOutput := string(buf[:n])

	if !strings.Contains(sseOutput, "bytesWritten") || !strings.Contains(sseOutput, "totalBytes") {
		t.Fatalf("expected SSE output to contain bytesWritten and totalBytes, got: %s", sseOutput)
	}
}
