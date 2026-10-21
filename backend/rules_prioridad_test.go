package main

import (
	"testing"
	"time"
)

func fechaUTC(anio int, mes time.Month, dia int) *time.Time {
	t := time.Date(anio, mes, dia, 0, 0, 0, 0, time.UTC)
	return &t
}

// Un caso por cada salida de PrioridadDeTarea, con los bordes de cada
// tramo (hoy, 7/8 días, 30/31 días) porque ahí es donde se equivoca un <=.
func TestPrioridadDeTarea(t *testing.T) {
	ahora := time.Date(2026, time.October, 20, 12, 0, 0, 0, time.UTC)

	casos := []struct {
		nombre   string
		tarea    Task
		esperada string
	}{
		{"archivada no tiene prioridad", Task{Status: "archivada", DueDate: fechaUTC(2026, time.October, 1)}, "ninguna"},
		{"completada no tiene prioridad", Task{Status: "completada", DueDate: fechaUTC(2026, time.October, 1)}, "ninguna"},
		{"pendiente sin fecha", Task{Status: "pendiente"}, "sin-fecha"},
		{"venció ayer", Task{Status: "pendiente", DueDate: fechaUTC(2026, time.October, 19)}, "vencida"},
		{"vence hoy", Task{Status: "pendiente", DueDate: fechaUTC(2026, time.October, 20)}, "hoy"},
		{"vence mañana", Task{Status: "pendiente", DueDate: fechaUTC(2026, time.October, 21)}, "esta-semana"},
		{"vence en 7 días (borde)", Task{Status: "pendiente", DueDate: fechaUTC(2026, time.October, 27)}, "esta-semana"},
		{"vence en 8 días (borde)", Task{Status: "pendiente", DueDate: fechaUTC(2026, time.October, 28)}, "este-mes"},
		{"vence en 30 días (borde)", Task{Status: "pendiente", DueDate: fechaUTC(2026, time.November, 19)}, "este-mes"},
		{"vence en 31 días", Task{Status: "pendiente", DueDate: fechaUTC(2026, time.November, 20)}, "mas-adelante"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := PrioridadDeTarea(c.tarea, ahora); got != c.esperada {
				t.Errorf("esperaba %q, dio %q", c.esperada, got)
			}
		})
	}
}

// La hora del día no cambia el resultado: se compara por fecha.
func TestPrioridadDeTarea_IgnoraLaHora(t *testing.T) {
	tarde := time.Date(2026, time.October, 20, 23, 59, 0, 0, time.UTC)
	tarea := Task{Status: "pendiente", DueDate: fechaUTC(2026, time.October, 20)}

	if got := PrioridadDeTarea(tarea, tarde); got != "hoy" {
		t.Errorf("esperaba \"hoy\", dio %q", got)
	}
}

func TestResumenDeLista_CuentaPorEstado(t *testing.T) {
	ahora := time.Date(2026, time.October, 20, 12, 0, 0, 0, time.UTC)
	tareas := []Task{
		{Status: "pendiente"},
		{Status: "pendiente", DueDate: fechaUTC(2026, time.October, 1)}, // vencida
		{Status: "completada"},
		{Status: "archivada"},
	}

	got := ResumenDeLista(tareas, ahora)
	esperado := Resumen{Total: 4, Pendientes: 2, Completadas: 1, Archivadas: 1, Vencidas: 1, PorcentajeHecho: 33}

	if got != esperado {
		t.Errorf("esperaba %+v, dio %+v", esperado, got)
	}
}

func TestResumenDeLista_SinTareasActivas(t *testing.T) {
	ahora := time.Date(2026, time.October, 20, 12, 0, 0, 0, time.UTC)

	casos := []struct {
		nombre string
		tareas []Task
		total  int
	}{
		{"lista vacía", nil, 0},
		{"todas archivadas", []Task{{Status: "archivada"}, {Status: "archivada"}}, 2},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got := ResumenDeLista(c.tareas, ahora)
			if got.Total != c.total || got.PorcentajeHecho != 0 {
				t.Errorf("esperaba total %d y 0%%, dio %+v", c.total, got)
			}
		})
	}
}

func TestResumenDeLista_TodoCompletado(t *testing.T) {
	ahora := time.Date(2026, time.October, 20, 12, 0, 0, 0, time.UTC)

	got := ResumenDeLista([]Task{{Status: "completada"}, {Status: "completada"}}, ahora)

	if got.PorcentajeHecho != 100 {
		t.Errorf("esperaba 100%%, dio %d%%", got.PorcentajeHecho)
	}
}

func TestResumenEtiqueta(t *testing.T) {
	casos := []struct {
		nombre   string
		resumen  Resumen
		esperada string
	}{
		{"lista vacía", Resumen{}, "lista vacía"},
		{"todo hecho", Resumen{Total: 2, Completadas: 2, PorcentajeHecho: 100}, "todo hecho"},
		{"con vencidas", Resumen{Total: 3, Pendientes: 2, Vencidas: 1, PorcentajeHecho: 33}, "hay tareas vencidas"},
		{"en progreso", Resumen{Total: 3, Pendientes: 2, PorcentajeHecho: 33}, "en progreso"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if got := c.resumen.Etiqueta(); got != c.esperada {
				t.Errorf("esperaba %q, dio %q", c.esperada, got)
			}
		})
	}
}
