package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func registerTaskRoutes(mux *http.ServeMux, s *Store) {
	// GET /api/lists/{listId}/tasks -> tareas de una lista
	mux.HandleFunc("GET /api/lists/{listId}/tasks", func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "falta identificar al usuario")
			return
		}

		listID, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}

		owner, err := s.ListOwner(r.Context(), listID)
		if err != nil {
			writeError(w, http.StatusNotFound, "lista no encontrada")
			return
		}
		if owner != userID {
			writeError(w, http.StatusForbidden, "no tenés acceso a esta lista")
			return
		}

		tasks, err := s.TasksByList(r.Context(), listID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		writeJSON(w, http.StatusOK, tasks)
	})

	// POST /api/lists/{listId}/tasks -> crear tarea
	mux.HandleFunc("POST /api/lists/{listId}/tasks", func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "falta identificar al usuario")
			return
		}

		listID, err := strconv.ParseInt(r.PathValue("listId"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}

		owner, err := s.ListOwner(r.Context(), listID)
		if err != nil {
			writeError(w, http.StatusNotFound, "lista no encontrada")
			return
		}
		if owner != userID {
			writeError(w, http.StatusForbidden, "no tenés acceso a esta lista")
			return
		}

		var body struct {
			Title   string  `json:"title"`
			DueDate *string `json:"due_date"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "cuerpo inválido")
			return
		}

		if err := ValidateTaskTitle(body.Title); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		var dueDate *time.Time
		if body.DueDate != nil && *body.DueDate != "" {
			parsed, err := time.Parse("2006-01-02", *body.DueDate)
			if err != nil {
				writeError(w, http.StatusBadRequest, "fecha inválida, formato esperado YYYY-MM-DD")
				return
			}
			dueDate = &parsed
		}

		if err := ValidateDueDate(dueDate, time.Now()); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		task, err := s.CreateTask(r.Context(), listID, body.Title, dueDate)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		writeJSON(w, http.StatusCreated, task)
	})

	// PATCH /api/tasks/{id} -> cambiar estado de una tarea
	mux.HandleFunc("PATCH /api/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "falta identificar al usuario")
			return
		}

		taskID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}

		task, err := s.TaskByID(r.Context(), taskID)
		if err != nil {
			writeError(w, http.StatusNotFound, "tarea no encontrada")
			return
		}

		owner, err := s.ListOwner(r.Context(), task.ListID)
		if err != nil || owner != userID {
			writeError(w, http.StatusForbidden, "no tenés acceso a esta tarea")
			return
		}

		var body struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "cuerpo inválido")
			return
		}

		if err := ValidateTransition(task.Status, body.Status); err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}

		if err := s.UpdateTaskStatus(r.Context(), taskID, body.Status); err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}

		task.Status = body.Status
		writeJSON(w, http.StatusOK, task)
	})

	// DELETE /api/tasks/{id} -> borrar tarea
	mux.HandleFunc("DELETE /api/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "falta identificar al usuario")
			return
		}

		taskID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "id inválido")
			return
		}

		task, err := s.TaskByID(r.Context(), taskID)
		if err != nil {
			writeError(w, http.StatusNotFound, "tarea no encontrada")
			return
		}

		owner, err := s.ListOwner(r.Context(), task.ListID)
		if err != nil || owner != userID {
			writeError(w, http.StatusForbidden, "no tenés acceso a esta tarea")
			return
		}

		if err := s.DeleteTask(r.Context(), taskID); err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
