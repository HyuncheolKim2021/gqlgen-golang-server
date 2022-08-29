package repository

import (
	"context"
	"gqlgen-golang-server/graph/model"
)

type TodoRepository struct {
}

func (tr *TodoRepository) GetTodos(ctx context.Context) []*model.Todo {
	user := &model.User{ID: "user-id-1", Name: "user"}
	todos := []*model.Todo{
		{ID: "1", User: user, Text: "Todo 1", Done: false},
		{ID: "2", User: user, Text: "Todo 2", Done: true},
		{ID: "3", User: user, Text: "Todo 3", Done: false},
	}
	return todos
}
