package service

import (
	"Bangseungjae/go-todo-app/entity"
	"Bangseungjae/go-todo-app/store"
	"context"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}
type TaskListener interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}
type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}
