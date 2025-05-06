package testhelper

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func ReadFixture(t *testing.T, filename string) []byte {
	t.Helper()
	dir, err := filepath.Abs("../fixture")
	if err != nil {
		t.Fatalf("Failed to resolve fixture directory. Check the path: %s", err)

	}
	bytes, err := os.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		t.Fatalf("Failed to read fixture. Check the path: %s", err)
	}
	return bytes
}

func NewMockAPI(t *testing.T, method string, filename string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(ReadFixture(t, filepath.Join(method, filename)))
		if err != nil {
			t.Fatalf("Failed to write fixture. Check the path: %s", err)
		}
	}))
}

func MarshalMockData[T any](t *testing.T, filename string) T {
	t.Helper()
	bytes := ReadFixture(t, filepath.Join("body", filename))

	var data T
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		t.Fatalf("error: %s", err.Error())
	}
	return data
}

func NewMockAPIWithStatus(t *testing.T, method string, filename string, status int, headers map[string]string, testRetry bool) *httptest.Server {
	t.Helper()
	requestCount := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}

		newStatus := status

		if testRetry && requestCount > 0 {
			newStatus = http.StatusOK
		}

		w.WriteHeader(newStatus)

		_, err := w.Write(ReadFixture(t, filepath.Join(method, filename)))
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
		requestCount++
	}))
}
