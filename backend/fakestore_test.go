package main

import (
	"context"
	"errors"
	"time"
)

// FakeStore implementa DataStore guardando todo en memoria (mapas), en vez
// de Postgres. Además cuenta llamadas a CreateList — eso es lo que convierte
// este doble de un simple Fake a algo que también sirve para verificar
// INTERACCIÓN (la definición de "mock" de la guía).
type FakeStore struct {
	users map[int64]*User
	lists map[int64]*List
	tasks map[int64]*Task

	nextUserID, nextListID, nextTaskID int64
	CreateListCalls                    int
}

func NewFakeStore() *FakeStore {
	return &FakeStore{users: map[int64]*User{}, lists: map[int64]*List{}, tasks: map[int64]*Task{}}
}

func (f *FakeStore) FindOrCreateUser(ctx context.Context, username string) (*User, error) {
	for _, u := range f.users {
		if u.Username == username {
			return u, nil
		}
	}
	f.nextUserID++
	u := &User{ID: f.nextUserID, Username: username}
	f.users[u.ID] = u
	return u, nil
}

func (f *FakeStore) CreateList(ctx context.Context, userID int64, name string) (*List, error) {
	f.CreateListCalls++
	f.nextListID++
	l := &List{ID: f.nextListID, UserID: userID, Name: name}
	f.lists[l.ID] = l
	return l, nil
}

func (f *FakeStore) ListsByUser(ctx context.Context, userID int64) ([]List, error) {
	result := []List{}
	for _, l := range f.lists {
		if l.UserID == userID {
			pending := 0
			for _, t := range f.tasks {
				if t.ListID == l.ID && t.Status == "pendiente" {
					pending++
				}
			}
			lc := *l
			lc.PendingCount = pending
			result = append(result, lc)
		}
	}
	return result, nil
}

func (f *FakeStore) ListOwner(ctx context.Context, listID int64) (int64, error) {
	l, ok := f.lists[listID]
	if !ok {
		return 0, errors.New("lista no encontrada")
	}
	return l.UserID, nil
}

func (f *FakeStore) PendingTaskCount(ctx context.Context, listID int64) (int, error) {
	count := 0
	for _, t := range f.tasks {
		if t.ListID == listID && t.Status == "pendiente" {
			count++
		}
	}
	return count, nil
}

func (f *FakeStore) DeleteList(ctx context.Context, listID int64) error {
	delete(f.lists, listID)
	return nil
}

func (f *FakeStore) CreateTask(ctx context.Context, listID int64, title string, dueDate *time.Time) (*Task, error) {
	f.nextTaskID++
	t := &Task{ID: f.nextTaskID, ListID: listID, Title: title, Status: "pendiente", DueDate: dueDate}
	f.tasks[t.ID] = t
	return t, nil
}

func (f *FakeStore) TasksByList(ctx context.Context, listID int64) ([]Task, error) {
	result := []Task{}
	for _, t := range f.tasks {
		if t.ListID == listID {
			result = append(result, *t)
		}
	}
	return result, nil
}

func (f *FakeStore) TaskByID(ctx context.Context, taskID int64) (*Task, error) {
	t, ok := f.tasks[taskID]
	if !ok {
		return nil, errors.New("tarea no encontrada")
	}
	return t, nil
}

func (f *FakeStore) UpdateTaskStatus(ctx context.Context, taskID int64, status string) error {
	t, ok := f.tasks[taskID]
	if !ok {
		return errors.New("tarea no encontrada")
	}
	t.Status = status
	return nil
}

func (f *FakeStore) DeleteTask(ctx context.Context, taskID int64) error {
	delete(f.tasks, taskID)
	return nil
}
