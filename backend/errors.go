package main

import "errors"

// Errores de reglas de negocio. Se definen una sola vez acá y se reutilizan
// tanto en las validaciones (rules.go) como en los handlers HTTP, para que
// el mensaje que ve el usuario y el que se testea en los unit tests sea
// siempre el mismo.
var (
	ErrEmptyTitle          = errors.New("el título no puede estar vacío")
	ErrTitleTooLong        = errors.New("el título no puede superar los 100 caracteres")
	ErrInvalidDueDate      = errors.New("la fecha de vencimiento no puede ser anterior a hoy")
	ErrInvalidTransition   = errors.New("transición de estado inválida")
	ErrListHasPendingTasks = errors.New("no se puede eliminar una lista con tareas pendientes")
)
