package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	todov1 "todo/gen/todo/v1"
	"todo/gen/todo/v1/todov1connect"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Todo struct {
	Id          int32
	TaskName    string
	Description string
	Done        bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Todo) toResponse() *todov1.Todo {
	return &todov1.Todo{
		Id:          t.Id,
		TaskName:    t.TaskName,
		Description: t.Description,
		Done:        t.Done,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}
}

type TodoServer struct {
	todos map[int32]Todo
	len   int32
}

func (s *TodoServer) AddTodo(
	ctx context.Context,
	req *connect.Request[todov1.AddTodoRequest],
) (*connect.Response[todov1.AddTodoResponse], error) {
	log.Println("Request headers: ", req.Header())

	timestamp := time.Now()
	s.len++
	id := s.len
	newTodo := Todo{
		Id:          id,
		TaskName:    fmt.Sprintf("%s", req.Msg.TaskName),
		Description: fmt.Sprintf("%s", req.Msg.Description),
		Done:        false,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}
	s.todos[id] = newTodo
	res := connect.NewResponse(&todov1.AddTodoResponse{Todo: newTodo.toResponse()})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}
func (s *TodoServer) DoneTodo(
	ctx context.Context,
	req *connect.Request[todov1.DoneTodoRequest],
) (*connect.Response[todov1.DoneTodoResponse], error) {
	log.Println("Request headers: ", req.Header())

	val, ok := s.todos[req.Msg.Id]
	if ok {
		val.Done = true
		s.todos[req.Msg.Id] = val
	}
	res := connect.NewResponse(&todov1.DoneTodoResponse{})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}
func (s *TodoServer) DeleteTodo(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {
	log.Println("Request headers: ", req.Header())
	delete(s.todos, req.Msg.Id)
	res := connect.NewResponse(&todov1.DeleteTodoResponse{})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}
func (s *TodoServer) ListTodo(
	ctx context.Context,
	req *connect.Request[todov1.ListTodoRequest],
) (*connect.Response[todov1.ListTodoResponse], error) {
	log.Println("Request headers: ", req.Header())
	fmt.Println(s.todos)
	len := len(s.todos)
	results := make([]*todov1.Todo, len)
	i := 0
	for _, val := range s.todos {
		results[i] = val.toResponse()
		i++
	}
	res := connect.NewResponse(&todov1.ListTodoResponse{Result: results})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}

func main() {
	todoServer := &TodoServer{
		todos: make(map[int32]Todo),
	}
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(todoServer)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{MaxHandlers: 100}),
	)
}
