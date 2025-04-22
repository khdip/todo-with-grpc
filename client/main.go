package main

import (
	"log"

	"google.golang.org/grpc"

	"practice/todo-with-grpc/client/todo"
)

func main() {
	conn, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}
	todoClient := todo.NewClient(conn)
	res, err := todoClient.GetTodo(3)
	if err != nil {
		log.Fatalf("Error while requesting: %s", err)
	}
	log.Printf("Response %#v", res)

	if err := todoClient.GetTodos(); err != nil {
		log.Fatalf("Error while requesting: %s", err)
	}

	if err := todoClient.SaveTodos(); err != nil {
		log.Fatalf("Error while saving: %s", err)
	}

	todoClient.BiDirectionalTodos()
}
