package main

import (
	"context"
	"net/http"
	"strconv"
	"testing"
)

func TestGetLists_SinHeader_NoAutorizado(t *testing.T) {
	mux := newTestMux(NewFakeStore())
	rec := doRequest(mux, http.MethodGet, "/api/lists", "", nil)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("esperaba 401, dio %d", rec.Code)
	}
}

func TestGetLists_ErrorDelStore_Devuelve500(t *testing.T) {
	mux := newTestMux(failingStore{NewFakeStore()})
	rec := doRequest(mux, http.MethodGet, "/api/lists", "1", nil)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("esperaba 500, dio %d", rec.Code)
	}
}

func TestPostLists_ErrorDelStore_Devuelve500(t *testing.T) {
	mux := newTestMux(failingStore{NewFakeStore()})
	rec := doRequest(mux, http.MethodPost, "/api/lists", "1", map[string]string{"name": "x"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("esperaba 500, dio %d", rec.Code)
	}
}

func TestDeleteList_ConTareasPendientes_Conflicto(t *testing.T) {
	store := NewFakeStore()
	ctx := context.Background()
	list, _ := store.CreateList(ctx, 1, "Lista")
	_, _ = store.CreateTask(ctx, list.ID, "pendiente", nil)
	mux := newTestMux(store)

	rec := doRequest(mux, http.MethodDelete, "/api/lists/"+strconv.FormatInt(list.ID, 10), "1", nil)
	if rec.Code != http.StatusConflict {
		t.Errorf("esperaba 409, dio %d", rec.Code)
	}
}

func TestDeleteList_ConTareasCompletadas_Permite(t *testing.T) {
	store := NewFakeStore()
	ctx := context.Background()
	list, _ := store.CreateList(ctx, 1, "Lista")
	task, _ := store.CreateTask(ctx, list.ID, "hecha", nil)
	_ = store.UpdateTaskStatus(ctx, task.ID, "completada")
	mux := newTestMux(store)

	rec := doRequest(mux, http.MethodDelete, "/api/lists/"+strconv.FormatInt(list.ID, 10), "1", nil)
	if rec.Code != http.StatusNoContent {
		t.Errorf("esperaba 204, dio %d", rec.Code)
	}
}

func TestDeleteList_IdInvalido_E_Inexistente(t *testing.T) {
	mux := newTestMux(NewFakeStore())

	if rec := doRequest(mux, http.MethodDelete, "/api/lists/abc", "1", nil); rec.Code != http.StatusBadRequest {
		t.Errorf("id inválido: esperaba 400, dio %d", rec.Code)
	}
	if rec := doRequest(mux, http.MethodDelete, "/api/lists/999", "1", nil); rec.Code != http.StatusNotFound {
		t.Errorf("inexistente: esperaba 404, dio %d", rec.Code)
	}
}
