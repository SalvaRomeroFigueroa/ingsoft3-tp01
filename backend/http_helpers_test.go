package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORS_Preflight_NoLlegaAlHandler(t *testing.T) {
	llamado := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { llamado = true })

	rec := httptest.NewRecorder()
	withCORS(next).ServeHTTP(rec, httptest.NewRequest(http.MethodOptions, "/api/lists", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("esperaba 200, dio %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("falta el header Access-Control-Allow-Origin")
	}
	if llamado {
		t.Error("un preflight OPTIONS no debería llegar al handler real")
	}
}

func TestWithCORS_RequestNormal_PasaAlHandler(t *testing.T) {
	llamado := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { llamado = true })

	rec := httptest.NewRecorder()
	withCORS(next).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/lists", nil))

	if !llamado {
		t.Error("un request normal debería llegar al handler")
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("falta el header Access-Control-Allow-Origin")
	}
}
