package todo

import (
	"context"
	"io"
	"log"
	tpb "practice/todo-with-grpc/proto/todo"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	tpb.UnimplementedTodoServiceServer
}

type Todo struct {
	ID          int64
	Title       string
	Description string
}

var todos = []Todo{
	{
		ID:          1,
		Title:       "This is title 1",
		Description: "This is description 1",
	},
	{
		ID:          2,
		Title:       "This is title 2",
		Description: "This is description 2",
	},
	{
		ID:          3,
		Title:       "This is title 3",
		Description: "This is description 3",
	},
}

func (s *Server) GetTodo(ctx context.Context, req *tpb.GetTodoRequest) (*tpb.GetTodoResponse, error) {
	log.Printf("Todo ID: %d", req.GetID())
	var todo Todo
	for _, t := range todos {
		if t.ID == req.GetID() {
			todo = t
			break
		}
	}

	if todo.ID == 0 {
		return &tpb.GetTodoResponse{}, status.Errorf(codes.NotFound, "Invalid ID")
	}

	return &tpb.GetTodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}, nil
}

func (s *Server) GetTodos(req *tpb.GetTodosRequest, stream tpb.TodoService_GetTodosServer) error {
	for _, t := range todos {
		err := stream.Send(&tpb.GetTodoResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
		})

		if err != nil {
			return status.Error(codes.Internal, "Failed to send.")
		}
		time.Sleep(time.Second * 2)
	}
	return nil
}

func (s *Server) SaveTodos(stream tpb.TodoService_SaveTodosServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&tpb.SaveTodoResponse{
				Body: "Streaming request ended.",
			})
		}

		if err != nil {
			return status.Error(codes.Unknown, err.Error())
		}

		log.Printf("%#v\n", req)
		todos = append(todos, Todo{
			ID:          req.Todo.ID,
			Title:       req.Todo.Title,
			Description: req.Todo.Description,
		})
		time.Sleep(time.Second * 2)
	}
}

func (s *Server) BiDirectionalTodos(stream tpb.TodoService_BiDirectionalTodosServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return status.Error(codes.Unknown, err.Error())
		}

		log.Printf("%#v\n", req)
		for _, todo := range todos {
			if todo.ID == req.ID {
				if err = stream.Send(&tpb.TodoResponse{
					Todo: &tpb.Todo{
						ID:          todo.ID,
						Title:       todo.Title,
						Description: todo.Description,
					},
				}); err != nil {
					return status.Error(codes.Unknown, err.Error())
				}
			}
		}
		time.Sleep(time.Second * 2)
	}
}
