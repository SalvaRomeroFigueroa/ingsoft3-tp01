package main

import (
	"context"
	"database/sql"
	"time"
)

// Store agrupa todo el acceso a la base de datos.
type Store struct {
	DB *sql.DB
}

// Migrate crea las tablas si no existen. Se corre al arrancar la app,
// así no hace falta un paso manual extra para levantar el esquema.
func (s *Store) Migrate(ctx context.Context) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS lists (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY,
	list_id INTEGER NOT NULL REFERENCES lists(id) ON DELETE CASCADE,
	title VARCHAR(100) NOT NULL,
	status VARCHAR(20) NOT NULL DEFAULT 'pendiente',
	due_date DATE
);
`
	_, err := s.DB.ExecContext(ctx, schema)
	return err
}

// FindOrCreateUser busca un usuario por nombre y lo crea si no existe.
func (s *Store) FindOrCreateUser(ctx context.Context, username string) (*User, error) {
	u := &User{}
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, username FROM users WHERE username = $1`, username,
	).Scan(&u.ID, &u.Username)

	if err == sql.ErrNoRows {
		err = s.DB.QueryRowContext(ctx,
			`INSERT INTO users (username) VALUES ($1) RETURNING id, username`, username,
		).Scan(&u.ID, &u.Username)
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// CreateList inserta una lista nueva para un usuario.
func (s *Store) CreateList(ctx context.Context, userID int64, name string) (*List, error) {
	l := &List{UserID: userID, Name: name}
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO lists (user_id, name) VALUES ($1, $2) RETURNING id`, userID, name,
	).Scan(&l.ID)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// ListsByUser devuelve las listas de un usuario junto con la cantidad
// de tareas pendientes de cada una (regla de CÁLCULO).
func (s *Store) ListsByUser(ctx context.Context, userID int64) ([]List, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT l.id, l.user_id, l.name,
			COALESCE(SUM(CASE WHEN t.status = 'pendiente' THEN 1 ELSE 0 END), 0) AS pending_count
		FROM lists l
		LEFT JOIN tasks t ON t.list_id = l.id
		WHERE l.user_id = $1
		GROUP BY l.id
		ORDER BY l.id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := []List{}
	for rows.Next() {
		var l List
		if err := rows.Scan(&l.ID, &l.UserID, &l.Name, &l.PendingCount); err != nil {
			return nil, err
		}
		lists = append(lists, l)
	}
	return lists, rows.Err()
}

// ListOwner devuelve el user_id dueño de una lista (para chequear autorización).
func (s *Store) ListOwner(ctx context.Context, listID int64) (int64, error) {
	var userID int64
	err := s.DB.QueryRowContext(ctx,
		`SELECT user_id FROM lists WHERE id = $1`, listID,
	).Scan(&userID)
	return userID, err
}

// PendingTaskCount cuenta las tareas pendientes de una lista.
func (s *Store) PendingTaskCount(ctx context.Context, listID int64) (int, error) {
	var count int
	err := s.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM tasks WHERE list_id = $1 AND status = 'pendiente'`, listID,
	).Scan(&count)
	return count, err
}

// DeleteList borra una lista (y en cascada sus tareas).
func (s *Store) DeleteList(ctx context.Context, listID int64) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM lists WHERE id = $1`, listID)
	return err
}

// CreateTask inserta una tarea nueva, siempre en estado "pendiente".
func (s *Store) CreateTask(ctx context.Context, listID int64, title string, dueDate *time.Time) (*Task, error) {
	t := &Task{ListID: listID, Title: title, Status: "pendiente", DueDate: dueDate}
	err := s.DB.QueryRowContext(ctx,
		`INSERT INTO tasks (list_id, title, status, due_date) VALUES ($1, $2, 'pendiente', $3) RETURNING id`,
		listID, title, dueDate,
	).Scan(&t.ID)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// TasksByList devuelve las tareas de una lista.
func (s *Store) TasksByList(ctx context.Context, listID int64) ([]Task, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, list_id, title, status, due_date FROM tasks WHERE list_id = $1 ORDER BY id`, listID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.ListID, &t.Title, &t.Status, &t.DueDate); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

// TaskByID busca una tarea por id.
func (s *Store) TaskByID(ctx context.Context, taskID int64) (*Task, error) {
	t := &Task{}
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, list_id, title, status, due_date FROM tasks WHERE id = $1`, taskID,
	).Scan(&t.ID, &t.ListID, &t.Title, &t.Status, &t.DueDate)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// UpdateTaskStatus actualiza el estado de una tarea.
func (s *Store) UpdateTaskStatus(ctx context.Context, taskID int64, status string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE tasks SET status = $1 WHERE id = $2`, status, taskID)
	return err
}

// DeleteTask borra una tarea.
func (s *Store) DeleteTask(ctx context.Context, taskID int64) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, taskID)
	return err
}
