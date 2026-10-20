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

// Resumen agrupa los contadores de una lista de tareas.
type Resumen struct {
	Total           int `json:"total"`
	Pendientes      int `json:"pendientes"`
	Completadas     int `json:"completadas"`
	Archivadas      int `json:"archivadas"`
	Vencidas        int `json:"vencidas"`
	PorcentajeHecho int `json:"porcentaje_hecho"`
}

// PrioridadDeTarea devuelve la prioridad de una tarea según su estado y su
// fecha de vencimiento. (Tiene varios caminos y, a propósito, ningún test.)
func PrioridadDeTarea(t Task, ahora time.Time) string {
	if t.Status == "archivada" || t.Status == "completada" {
		return "ninguna"
	}
	if t.DueDate == nil {
		return "sin-fecha"
	}
	hoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, time.UTC)
	vence := time.Date(t.DueDate.Year(), t.DueDate.Month(), t.DueDate.Day(), 0, 0, 0, 0, time.UTC)
	dias := int(vence.Sub(hoy).Hours() / 24)
	switch {
	case dias < 0:
		return "vencida"
	case dias == 0:
		return "hoy"
	case dias <= 7:
		return "esta-semana"
	case dias <= 30:
		return "este-mes"
	default:
		return "mas-adelante"
	}
}

// ResumenDeLista cuenta las tareas por estado y calcula el avance sobre las
// tareas activas (las archivadas no cuentan).
func ResumenDeLista(tareas []Task, ahora time.Time) Resumen {
	var r Resumen
	for _, t := range tareas {
		r.Total++
		switch t.Status {
		case "pendiente":
			r.Pendientes++
			if PrioridadDeTarea(t, ahora) == "vencida" {
				r.Vencidas++
			}
		case "completada":
			r.Completadas++
		case "archivada":
			r.Archivadas++
		}
	}
	activas := r.Total - r.Archivadas
	if activas == 0 {
		return r
	}
	r.PorcentajeHecho = r.Completadas * 100 / activas
	return r
}

// Etiqueta resume el estado de la lista en una frase corta.
func (r Resumen) Etiqueta() string {
	if r.Total == 0 {
		return "lista vacía"
	}
	if r.PorcentajeHecho == 100 {
		return "todo hecho"
	}
	if r.Vencidas > 0 {
		return "hay tareas vencidas"
	}
	return "en progreso"
}
