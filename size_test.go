package bodysize_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	bodysize "github.com/Russia9/body-size-limit"
)

func TestSuccess(t *testing.T) {
	cfg := bodysize.CreateConfig()
	cfg.Limit = 5

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := bodysize.New(ctx, next, cfg, "body-size")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", bytes.NewReader([]byte{1, 2, 3}))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertStatus(t, recorder, 200)
}

func TestError(t *testing.T) {
	cfg := bodysize.CreateConfig()
	cfg.Limit = 5

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := bodysize.New(ctx, next, cfg, "body-size")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", bytes.NewReader([]byte{1, 2, 3, 4, 5, 6}))
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertStatus(t, recorder, 413)
}

func assertStatus(t *testing.T, resp *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if resp.Code != expected {
		t.Errorf("invalid status code: %d", resp.Code)
	}
}
