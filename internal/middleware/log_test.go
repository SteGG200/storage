package middleware_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SteGG200/storage/internal/middleware"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	// Save the original logger flags and output
	originalFlags := log.Flags()
	originalOutput := log.Writer()
	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	log.SetOutput(&buf)
	log.SetFlags(0) // Remove date/time prefix for easier string matching

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	loggedHandler := middleware.Logger(handler)

	req := httptest.NewRequest("GET", "/test-path", nil)
	rr := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rr.Code)
	}

	output := strings.TrimSpace(buf.String())
	// In development build (!production tag), it should log.
	// In production build, it is a no-op so output is empty.
	if output != "" {
		if !strings.HasPrefix(output, "GET /test-path") {
			t.Errorf("unexpected log output: %q", output)
		}
		if !strings.HasSuffix(output, "ms") {
			t.Errorf("expected log output to end with 'ms', got: %q", output)
		}
	}
}
