package main

import (
	"strings"
	"time"
)

// MaxTitleLength es el largo máximo permitido para el título de una tarea.
const MaxTitleLength = 100

// ValidateTaskTitle: regla de VALIDACIÓN.
// El título no puede estar vacío ni superar los 100 caracteres.
func ValidateTaskTitle(title string) error {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return ErrEmptyTitle
	}
	if len(trimmed) > MaxTitleLength {
		return ErrTitleTooLong
	}
	return nil
}

// ValidateDueDate: regla de VALIDACIÓN.
// La fecha de vencimiento (si se especifica) no puede ser anterior a hoy.
func ValidateDueDate(due *time.Time, now time.Time) error {
	if due == nil {
		return nil
	}
	y1, m1, d1 := due.Date()
	y2, m2, d2 := now.Date()
	dueDay := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	today := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	if dueDay.Before(today) {
		return ErrInvalidDueDate
	}
	return nil
}

// allowedTransitions: regla de TRANSICIÓN DE ESTADO.
// pendiente <-> completada, y ambas pueden pasar a archivada.
// archivada es un estado terminal: no puede volver atrás.
var allowedTransitions = map[string]map[string]bool{
	"pendiente":  {"completada": true, "archivada": true},
	"completada": {"pendiente": true, "archivada": true},
	"archivada":  {},
}

// ValidateTransition valida si se puede pasar de un estado a otro.
func ValidateTransition(current, next string) error {
	if current == next {
		return nil
	}
	allowed, ok := allowedTransitions[current]
	if !ok || !allowed[next] {
		return ErrInvalidTransition
	}
	return nil
}

// CanDeleteList: regla de RESTRICCIÓN.
// No se puede eliminar una lista que todavía tiene tareas pendientes.
func CanDeleteList(pendingCount int) error {
	if pendingCount > 0 {
		return ErrListHasPendingTasks
	}
	return nil
}
