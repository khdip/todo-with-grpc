package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	tpb "practice/todo-with-grpc/proto/todo"
	"practice/todo-with-grpc/server/todo"
)

func main() {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}

	grpcServer := grpc.NewServer()
	s := todo.Server{}
	tpb.RegisterTodoServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
