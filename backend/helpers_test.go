package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
)

func newTestMux(store DataStore) *http.ServeMux {
	mux := http.NewServeMux()
	registerUserRoutes(mux, store)
	registerListRoutes(mux, store)
	registerTaskRoutes(mux, store)
	return mux
}

func doRequest(mux http.Handler, method, path, userID string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	if userID != "" {
		req.Header.Set("X-User-Id", userID)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// failingStore es un FakeStore que falla en ListsByUser y CreateList,
// para probar cómo reacciona el handler cuando la base se cae.
type failingStore struct{ *FakeStore }

func (failingStore) ListsByUser(ctx context.Context, userID int64) ([]List, error) {
	return nil, errors.New("fallo simulado de la base")
}

func (failingStore) CreateList(ctx context.Context, userID int64, name string) (*List, error) {
	return nil, errors.New("fallo simulado de la base")
}
