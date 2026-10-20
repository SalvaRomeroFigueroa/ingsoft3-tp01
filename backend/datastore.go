package main

import (
	"context"
	"time"
)

// DataStore es lo que un handler necesita, sin saber si del otro lado hay
// Postgres real o un doble de test.
type DataStore interface {
	FindOrCreateUser(ctx context.Context, username string) (*User, error)

	CreateList(ctx context.Context, userID int64, name string) (*List, error)
	ListsByUser(ctx context.Context, userID int64) ([]List, error)
	ListOwner(ctx context.Context, listID int64) (int64, error)
	PendingTaskCount(ctx context.Context, listID int64) (int, error)
	DeleteList(ctx context.Context, listID int64) error

	CreateTask(ctx context.Context, listID int64, title string, dueDate *time.Time) (*Task, error)
	TasksByList(ctx context.Context, listID int64) ([]Task, error)
	TaskByID(ctx context.Context, taskID int64) (*Task, error)
	UpdateTaskStatus(ctx context.Context, taskID int64, status string) error
	DeleteTask(ctx context.Context, taskID int64) error
}
