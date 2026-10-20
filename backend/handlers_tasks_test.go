package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func tasksPath(listID int64) string {
	return "/api/lists/" + strconv.FormatInt(listID, 10) + "/tasks"
}

func taskPath(taskID int64) string {
	return "/api/tasks/" + strconv.FormatInt(taskID, 10)
}

func TestCreateTask_Validaciones(t *testing.T) {
	pasada := "2000-01-01"
	futura := "2999-01-01"
	formatoMalo := "31/12/2030"

	cases := []struct {
		nombre   string
		title    string
		dueDate  *string
		esperado int
	}{
		{"título válido sin fecha", "Comprar pan", nil, http.StatusCreated},
		{"título válido con fecha futura", "Comprar pan", &futura, http.StatusCreated},
		{"título de exactamente 100 caracteres", strings.Repeat("a", 100), nil, http.StatusCreated},
		{"título vacío", "   ", nil, http.StatusBadRequest},
		{"título de 101 caracteres", strings.Repeat("a", 101), nil, http.StatusBadRequest},
		{"fecha en el pasado", "Comprar pan", &pasada, http.StatusBadRequest},
		{"fecha con formato inválido", "Comprar pan", &formatoMalo, http.StatusBadRequest},
	}

	for _, c := range cases {
		t.Run(c.nombre, func(t *testing.T) {
			store := NewFakeStore()
			list, _ := store.CreateList(context.Background(), 1, "Lista")
			mux := newTestMux(store)

			body := map[string]interface{}{"title": c.title, "due_date": c.dueDate}
			rec := doRequest(mux, http.MethodPost, tasksPath(list.ID), "1", body)

			if rec.Code != c.esperado {
				t.Errorf("esperaba %d, dio %d: %s", c.esperado, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestCreateTask_Autorizacion(t *testing.T) {
	store := NewFakeStore()
	list, _ := store.CreateList(context.Background(), 1, "Lista de 1")
	mux := newTestMux(store)
	body := map[string]interface{}{"title": "algo"}

	if rec := doRequest(mux, http.MethodPost, tasksPath(list.ID), "2", body); rec.Code != http.StatusForbidden {
		t.Errorf("lista ajena: esperaba 403, dio %d", rec.Code)
	}
	if rec := doRequest(mux, http.MethodPost, tasksPath(999), "1", body); rec.Code != http.StatusNotFound {
		t.Errorf("lista inexistente: esperaba 404, dio %d", rec.Code)
	}
	if rec := doRequest(mux, http.MethodPost, tasksPath(list.ID), "", body); rec.Code != http.StatusUnauthorized {
		t.Errorf("sin usuario: esperaba 401, dio %d", rec.Code)
	}
}

func TestUpdateTaskStatus_Transiciones(t *testing.T) {
	cases := []struct {
		nombre   string
		inicial  string
		nuevo    string
		esperado int
	}{
		{"pendiente a completada", "pendiente", "completada", http.StatusOK},
		{"completada a pendiente", "completada", "pendiente", http.StatusOK},
		{"pendiente a archivada", "pendiente", "archivada", http.StatusOK},
		{"archivada no vuelve a pendiente", "archivada", "pendiente", http.StatusConflict},
		{"archivada no pasa a completada", "archivada", "completada", http.StatusConflict},
		{"estado inexistente", "pendiente", "inventado", http.StatusConflict},
	}

	for _, c := range cases {
		t.Run(c.nombre, func(t *testing.T) {
			store := NewFakeStore()
			ctx := context.Background()
			list, _ := store.CreateList(ctx, 1, "Lista")
			task, _ := store.CreateTask(ctx, list.ID, "tarea", nil)
			_ = store.UpdateTaskStatus(ctx, task.ID, c.inicial)
			mux := newTestMux(store)

			rec := doRequest(mux, http.MethodPatch, taskPath(task.ID), "1", map[string]string{"status": c.nuevo})

			if rec.Code != c.esperado {
				t.Errorf("esperaba %d, dio %d: %s", c.esperado, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestUpdateTaskStatus_Autorizacion(t *testing.T) {
	store := NewFakeStore()
	ctx := context.Background()
	list, _ := store.CreateList(ctx, 1, "Lista de 1")
	task, _ := store.CreateTask(ctx, list.ID, "tarea", nil)
	mux := newTestMux(store)
	body := map[string]string{"status": "completada"}

	if rec := doRequest(mux, http.MethodPatch, taskPath(task.ID), "2", body); rec.Code != http.StatusForbidden {
		t.Errorf("tarea ajena: esperaba 403, dio %d", rec.Code)
	}
	if rec := doRequest(mux, http.MethodPatch, taskPath(999), "1", body); rec.Code != http.StatusNotFound {
		t.Errorf("tarea inexistente: esperaba 404, dio %d", rec.Code)
	}
	if rec := doRequest(mux, http.MethodPatch, taskPath(task.ID), "", body); rec.Code != http.StatusUnauthorized {
		t.Errorf("sin usuario: esperaba 401, dio %d", rec.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	store := NewFakeStore()
	ctx := context.Background()
	list, _ := store.CreateList(ctx, 1, "Lista de 1")
	task, _ := store.CreateTask(ctx, list.ID, "tarea", nil)
	mux := newTestMux(store)

	if rec := doRequest(mux, http.MethodDelete, taskPath(task.ID), "2", nil); rec.Code != http.StatusForbidden {
		t.Errorf("ajeno: esperaba 403, dio %d", rec.Code)
	}
	if rec := doRequest(mux, http.MethodDelete, taskPath(999), "1", nil); rec.Code != http.StatusNotFound {
		t.Errorf("inexistente: esperaba 404, dio %d", rec.Code)
	}
	if rec := doRequest(mux, http.MethodDelete, taskPath(task.ID), "1", nil); rec.Code != http.StatusNoContent {
		t.Errorf("dueño: esperaba 204, dio %d", rec.Code)
	}
	if _, err := store.TaskByID(ctx, task.ID); err == nil {
		t.Error("la tarea debería haberse borrado del store")
	}
}

func TestGetTasks(t *testing.T) {
	store := NewFakeStore()
	ctx := context.Background()
	list, _ := store.CreateList(ctx, 1, "Lista de 1")
	_, _ = store.CreateTask(ctx, list.ID, "tarea", nil)
	mux := newTestMux(store)

	rec := doRequest(mux, http.MethodGet, tasksPath(list.ID), "1", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("dueño: esperaba 200, dio %d", rec.Code)
	}
	var tasks []Task
	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("respuesta no es JSON válido: %v", err)
	}
	if len(tasks) != 1 {
		t.Errorf("esperaba 1 tarea, dio %d", len(tasks))
	}

	if rec := doRequest(mux, http.MethodGet, tasksPath(list.ID), "2", nil); rec.Code != http.StatusForbidden {
		t.Errorf("ajeno: esperaba 403, dio %d", rec.Code)
	}
}
