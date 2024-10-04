package service

import (
	"Bangseungjae/go-todo-app/entity"
	"Bangseungjae/go-todo-app/store"
	"context"
	"fmt"
)

type ListTask struct {
	DB   store.Queryer
	Repo TaskListener
}

func (l *ListTask) ListTasks(ctx context.Context) (entity.Tasks, error) {
	ts, err := l.Repo.ListTasks(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return ts, nil
}
