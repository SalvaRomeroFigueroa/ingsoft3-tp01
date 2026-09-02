package main

import "time"

// User representa a un usuario de la aplicación.
// No hay contraseña: para mantener el TP simple, el "login" es solo
// identificarse por nombre de usuario (se crea si no existe).
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// List representa una lista de tareas, propiedad de un usuario.
type List struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	Name         string `json:"name"`
	PendingCount int    `json:"pending_count"`
}

// Task representa una tarea dentro de una lista.
// Status puede ser: "pendiente", "completada" o "archivada".
type Task struct {
	ID      int64      `json:"id"`
	ListID  int64      `json:"list_id"`
	Title   string     `json:"title"`
	Status  string     `json:"status"`
	DueDate *time.Time `json:"due_date,omitempty"`
}
