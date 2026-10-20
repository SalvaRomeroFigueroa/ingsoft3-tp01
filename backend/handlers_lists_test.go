package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// CON DOBLE: inyecta FakeStore en vez de Postgres, y el assert verifica
// tanto el resultado HTTP como la INTERACCIÓN (que CreateList se llamó
// exactamente una vez) — eso es lo que lo hace un mock, no solo un stub.
func TestPostLists_CreaListaConDoble(t *testing.T) {
	store := NewFakeStore()
	mux := http.NewServeMux()
	registerListRoutes(mux, store)

	body, _ := json.Marshal(map[string]string{"name": "Compras"})
	req := httptest.NewRequest(http.MethodPost, "/api/lists", bytes.NewReader(body))
	req.Header.Set("X-User-Id", "1")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, dio %d: %s", rec.Code, rec.Body.String())
	}
	if store.CreateListCalls != 1 {
		t.Errorf("esperaba 1 llamada a CreateList, hubo %d", store.CreateListCalls)
	}

	var got List
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}
	if got.Name != "Compras" {
		t.Errorf("esperaba nombre 'Compras', dio %q", got.Name)
	}
}

// CASO DE ERROR: nombre vacío rechazado con 400, sin tocar el store.
func TestPostLists_NombreVacio_Rechazado(t *testing.T) {
	store := NewFakeStore()
	mux := http.NewServeMux()
	registerListRoutes(mux, store)

	body, _ := json.Marshal(map[string]string{"name": "   "})
	req := httptest.NewRequest(http.MethodPost, "/api/lists", bytes.NewReader(body))
	req.Header.Set("X-User-Id", "1")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, dio %d", rec.Code)
	}
	if store.CreateListCalls != 0 {
		t.Errorf("no debería haber llamado a CreateList, se llamó %d veces", store.CreateListCalls)
	}
}

// PARAMETRIZADO (tabla), el equivalente en Go de [Theory]+[InlineData]:
// varios casos de autorización sin duplicar el test.
func TestDeleteList_Autorizacion(t *testing.T) {
	cases := []struct {
		nombre         string
		userHeader     string
		listOwnerID    int64
		statusEsperado int
	}{
		{"dueño borra su propia lista", "1", 1, http.StatusNoContent},
		{"otro usuario no puede borrarla", "2", 1, http.StatusForbidden},
	}

	for _, c := range cases {
		t.Run(c.nombre, func(t *testing.T) {
			store := NewFakeStore()
			list, _ := store.CreateList(context.Background(), c.listOwnerID, "Lista de prueba")

			mux := http.NewServeMux()
			registerListRoutes(mux, store)

			req := httptest.NewRequest(http.MethodDelete, "/api/lists/"+strconv.FormatInt(list.ID, 10), nil)
			req.Header.Set("X-User-Id", c.userHeader)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != c.statusEsperado {
				t.Errorf("esperaba %d, dio %d", c.statusEsperado, rec.Code)
			}
		})
	}
}
