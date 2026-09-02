package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func registerUserRoutes(mux *http.ServeMux, s *Store) {
	// POST /api/users/login
	// "Login" simplificado: si el username no existe se crea. No hay
	// contraseña porque no es el foco de este TP (foco: contenedores).
	mux.HandleFunc("POST /api/users/login", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string `json:"username"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeError(w, http.StatusBadRequest, "cuerpo inválido")
			return
		}

		username := strings.TrimSpace(body.Username)
		if username == "" {
			writeError(w, http.StatusBadRequest, "el nombre de usuario no puede estar vacío")
			return
		}

		user, err := s.FindOrCreateUser(r.Context(), username)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "error interno")
			return
		}
		writeJSON(w, http.StatusOK, user)
	})
}
