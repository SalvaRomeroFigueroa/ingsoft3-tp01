package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func registerListRoutes(mux *http.ServeMux, s *Store) {
	// GET /api/lists -> listas del usuario autenticado
	mux.HandleFunc("GET /api/lists", func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "falta identificar al usuario")
			return
		}
		lists, err := s.ListsByUser(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		writeJSON(w, http.StatusOK, lists)
	})

	// POST /api/lists -> crear lista
	mux.HandleFunc("POST /api/lists", func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "falta identificar al usuario")
			return
		}

		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "cuerpo inválido")
			return
		}

		name := strings.TrimSpace(body.Name)
		if name == "" {
			writeError(w, http.StatusBadRequest, "el nombre de la lista no puede estar vacío")
			return
		}

		list, err := s.CreateList(r.Context(), userID, name)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		writeJSON(w, http.StatusCreated, list)
	})

	// DELETE /api/lists/{id} -> borrar lista (si no tiene tareas pendientes)
	mux.HandleFunc("DELETE /api/lists/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "falta identificar al usuario")
			return
		}

		listID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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

		pending, err := s.PendingTaskCount(r.Context(), listID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		if err := CanDeleteList(pending); err != nil {
			writeError(w, http.StatusConflict, err.Error())
			return
		}

		if err := s.DeleteList(r.Context(), listID); err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
