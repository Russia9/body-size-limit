package body_size_limit_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
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
	assertBody(t, req, []byte{1, 2, 3})
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
	assertBody(t, req, []byte{})
}

func assertStatus(t *testing.T, resp *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if resp.Code != expected {
		t.Errorf("invalid status code: %d", resp.Code)
	}
}

func assertBody(t *testing.T, req *http.Request, expected []byte) {
	t.Helper()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Errorf("Error while reading body: %e", err)
	}

	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Unexpected body: %v", body)
	}
}
