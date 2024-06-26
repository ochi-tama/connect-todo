package main

import (
	"context"
	"log"
	"net/http"
	todov1 "todo/gen/todo/v1"
	"todo/gen/todo/v1/todov1connect"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type TodoServer struct{}

func (s *TodoServer) AddTodo(
	ctx context.Context,
	req *connect.Request[todov1.AddTodoRequest],
) (*connect.Response[todov1.AddTodoResponse], error) {
	log.Println("Request headers: ", req.Header())

	todo := todov1.Todo{}
	res := connect.NewResponse(&todov1.AddTodoResponse{Todo: &todo})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}
func (s *TodoServer) DoneTodo(
	ctx context.Context,
	req *connect.Request[todov1.DoneTodoRequest],
) (*connect.Response[todov1.DoneTodoResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&todov1.DoneTodoResponse{})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}
func (s *TodoServer) DeleteTodo(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&todov1.DeleteTodoResponse{})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}
func (s *TodoServer) ListTodo(
	ctx context.Context,
	req *connect.Request[todov1.ListTodoRequest],
) (*connect.Response[todov1.ListTodoResponse], error) {
	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&todov1.ListTodoResponse{})
	res.Header().Set("Todo-Version", "v1")
	return res, nil
}

func main() {
	todoServer := &TodoServer{}
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(todoServer)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{MaxHandlers: 100}),
	)
}
