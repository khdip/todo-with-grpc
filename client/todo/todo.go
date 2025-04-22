package todo

import (
	"context"
	"io"
	"log"
	tpb "practice/todo-with-grpc/proto/todo"
	"time"

	"google.golang.org/grpc"
)

type Client struct {
	client tpb.TodoServiceClient
}

func NewClient(conn grpc.ClientConnInterface) Client {
	return Client{
		client: tpb.NewTodoServiceClient(conn),
	}
}

func (c *Client) GetTodo(id int64) (*tpb.GetTodoResponse, error) {
	return c.client.GetTodo(context.Background(), &tpb.GetTodoRequest{
		ID: id,
	})
}

func (c *Client) GetTodos() error {
	stream, err := c.client.GetTodos(context.Background(), &tpb.GetTodosRequest{})
	if err != nil {
		return err
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("Received: %#v", res)
	}
}

type Todo struct {
	ID          int64
	Title       string
	Description string
}

var todos = []Todo{
	{
		ID:          4,
		Title:       "This is title 4",
		Description: "This is description 4",
	},
	{
		ID:          5,
		Title:       "This is title 5",
		Description: "This is description 5",
	},
	{
		ID:          6,
		Title:       "This is title 6",
		Description: "This is description 6",
	},
}

func (c *Client) SaveTodos() error {
	stream, err := c.client.SaveTodos(context.Background())
	if err != nil {
		return err
	}

	for _, todo := range todos {
		if err := stream.Send(&tpb.SaveTodoRequest{
			Todo: &tpb.Todo{
				ID:          todo.ID,
				Title:       todo.Title,
				Description: todo.Description,
			},
		}); err != nil {
			return err
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil
	}
	log.Println("Response from server: ", res.Body)
	return nil
}

func (c *Client) BiDirectionalTodos() {
	stream, err := c.client.BiDirectionalTodos(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}

	for _, id := range []int64{1, 2, 3} {
		if err := stream.Send(&tpb.TodoRequest{
			ID: id,
		}); err != nil {
			log.Fatalf("Failed to call: %s", err)
		}
		time.Sleep(time.Second * 2)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(waitc)
					return
				}
				log.Fatal(err)
			}
			log.Printf("Response from bidirectional server: %#v\n", res)
		}
	}()

	stream.CloseSend()
	<-waitc
}
