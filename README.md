# gqlgenで始めるGoのGraphQLサーバー作り（初級編）

```go
$ go mod init <PROJECT_NAME>
$ touch tools.go
```

### gqlgenを導入する

```go
//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
)
```
```go
$ go mod tidy
```
`tools.go`はいらないので削除しても良い

### gqlgenをinitialize

```go
$ go run github.com/99designs/gqlgen init
```

`/graph`や`gqlgen.yml`, `server.go`が生成される

`/graph`の中身
```go
$ ls -R graph
generated           model               resolver.go         schema.graphqls     schema.resolvers.go

graph/generated:
generated.go

graph/model:
models_gen.go
```

簡単なTODOを表現するスキーマができている
```graphql
# schema.graphqls
# GraphQL schema example
#
# https://gqlgen.com/getting-started/

type Todo {
id: ID!
text: String!
done: Boolean!
user: User!
}

type User {
id: ID!
name: String!
}

type Query {
todos: [Todo!]!
}

input NewTodo {
text: String!
userId: String!
}

type Mutation {
createTodo(input: NewTodo!): Todo!
}
```

### gqlgenが生成した`schema.resolvers.go`
```go
package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"gqlgen-golang-server/graph/generated"
	"gqlgen-golang-server/graph/model"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

```
```go
// 省略
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}
// 省略
```
ここを実装して返せばおk。

ますはRepositoryの実装
```go
// $ mkdir repository && touch todo_repository.go
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
```

今回はDBとの繋ぎはやらないので仮データをダミーデータを返しましょう

Service層の実装
```go
// $ mkdir service && touch todo_service.go
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
```

serviceのビルダ関数とRepositoryを呼び、ダミーデータを取得します

`/graph/resolver.go`からサービスを呼べるようにし、
```go
func (r *Resolver) NewTodoService() *service.TodoService {
	repository := repository.TodoRepository{}
	return service.NewTodoService(repository)
}
```

`/graph/schema.resolvers.go`ではこの関数よ呼び、Todosを返しましょう
```go
// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	return r.Resolver.NewTodoService().GetTodos(ctx), nil
}
```

これでQuery todosの完成です！実際Queryを実行してみましょう。
```go
$ go run server.go
```

`localhost:8080`にアクセスし、以下のクエリを投げます
```graphql
query {
  todos {
    id
    user{
      id
      name
    }
    text
    done
  }
}
```
そうすると、
```json
{
  "data": {
    "todos": [
      {
        "id": "1",
        "user": {
          "id": "user-id-1",
          "name": "user"
        },
        "text": "Todo 1",
        "done": false
      },
      {
        "id": "2",
        "user": {
          "id": "user-id-1",
          "name": "user"
        },
        "text": "Todo 2",
        "done": true
      },
      {
        "id": "3",
        "user": {
          "id": "user-id-1",
          "name": "user"
        },
        "text": "Todo 3",
        "done": false
      }
    ]
  }
}
```
のように、ダミーデータがちゃんと返ってくることが分かります！