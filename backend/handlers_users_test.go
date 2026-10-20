package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin_CreaYReutilizaUsuario(t *testing.T) {
	mux := newTestMux(NewFakeStore())
	body := map[string]string{"username": "salva"}

	rec := doRequest(mux, http.MethodPost, "/api/users/login", "", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("esperaba 200, dio %d", rec.Code)
	}
	var u1 User
	if err := json.Unmarshal(rec.Body.Bytes(), &u1); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}
	if u1.Username != "salva" || u1.ID == 0 {
		t.Errorf("usuario inesperado: %+v", u1)
	}

	rec = doRequest(mux, http.MethodPost, "/api/users/login", "", body)
	var u2 User
	_ = json.Unmarshal(rec.Body.Bytes(), &u2)
	if u2.ID != u1.ID {
		t.Errorf("el mismo username debería devolver el mismo id: %d vs %d", u1.ID, u2.ID)
	}
}

func TestLogin_Invalidos(t *testing.T) {
	mux := newTestMux(NewFakeStore())

	rec := doRequest(mux, http.MethodPost, "/api/users/login", "", map[string]string{"username": "   "})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("username vacío: esperaba 400, dio %d", rec.Code)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader("esto no es json"))
	recRaw := httptest.NewRecorder()
	mux.ServeHTTP(recRaw, req)
	if recRaw.Code != http.StatusBadRequest {
		t.Errorf("cuerpo no JSON: esperaba 400, dio %d", recRaw.Code)
	}
}
