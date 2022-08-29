package graph

import (
	"gqlgen-golang-server/repository"
	"gqlgen-golang-server/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

func (r *Resolver) NewTodoService() *service.TodoService {
	repository := repository.TodoRepository{}
	return service.NewTodoService(repository)
}
