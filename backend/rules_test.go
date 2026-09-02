package main

import (
	"testing"
	"time"
)

// Estos son solo algunos tests de ejemplo, a modo de muestra de cómo
// testear las reglas de rules.go sin necesidad de base de datos.
// En el TP5 hay que llegar a 8 tests de backend: con las 4 reglas de
// este archivo (título, fecha, transición, borrado de lista) alcanza
// agregando los casos borde que falten.

func TestValidateTaskTitle(t *testing.T) {
	if err := ValidateTaskTitle("Comprar pan"); err != nil {
		t.Errorf("esperaba título válido, dio error: %v", err)
	}
	if err := ValidateTaskTitle("   "); err != ErrEmptyTitle {
		t.Errorf("esperaba ErrEmptyTitle, dio: %v", err)
	}
}

func TestValidateDueDate(t *testing.T) {
	now := time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)
	ayer := now.AddDate(0, 0, -1)
	if err := ValidateDueDate(&ayer, now); err != ErrInvalidDueDate {
		t.Errorf("esperaba ErrInvalidDueDate, dio: %v", err)
	}
	if err := ValidateDueDate(nil, now); err != nil {
		t.Errorf("esperaba nil para fecha no especificada, dio: %v", err)
	}
}

func TestValidateTransition(t *testing.T) {
	if err := ValidateTransition("pendiente", "completada"); err != nil {
		t.Errorf("transición permitida, dio error: %v", err)
	}
	if err := ValidateTransition("archivada", "pendiente"); err != ErrInvalidTransition {
		t.Errorf("esperaba ErrInvalidTransition, dio: %v", err)
	}
}

func TestCanDeleteList(t *testing.T) {
	if err := CanDeleteList(0); err != nil {
		t.Errorf("sin pendientes debería poder borrarse, dio: %v", err)
	}
	if err := CanDeleteList(2); err != ErrListHasPendingTasks {
		t.Errorf("esperaba ErrListHasPendingTasks, dio: %v", err)
	}
}
