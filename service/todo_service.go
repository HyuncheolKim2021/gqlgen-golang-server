package service

import (
	"context"
	"gqlgen-golang-server/graph/model"
	"gqlgen-golang-server/repository"
)

type TodoService struct {
	todoRepository repository.TodoRepository
}

func NewTodoService(todoRepository repository.TodoRepository) *TodoService {
	return &TodoService{todoRepository: todoRepository}
}

func (ts *TodoService) GetTodos(ctx context.Context) []*model.Todo {
	return ts.todoRepository.GetTodos(ctx)
}
