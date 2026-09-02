package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// writeJSON serializa v como JSON con el status code indicado.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError responde un error en formato JSON: {"error": "..."}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// withCORS habilita CORS para que el frontend (servido en otro origen
// durante el desarrollo local) pueda llamar a la API sin problemas.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-Id")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// userIDFromRequest lee el usuario "logueado" desde el header X-User-Id.
// No es un sistema de autenticación real (no hace falta para este TP),
// pero alcanza para poder testear la regla de autorización.
func userIDFromRequest(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.Header.Get("X-User-Id"), 10, 64)
}
