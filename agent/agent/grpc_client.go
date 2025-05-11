package agent

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/ZolotarevAlexandr/yl_final/grpc"
)

// Client wraps the gRPC client connection and methods
type Client struct {
	cli  pb.TaskServiceClient
	conn *grpc.ClientConn
}

// NewClient creates a new gRPC client connection to the orchestrator
func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		cli:  pb.NewTaskServiceClient(conn),
		conn: conn,
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// FetchTask retrieves a task from the orchestrator
func (c *Client) FetchTask() (*pb.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.cli.GetTask(ctx, &emptypb.Empty{})
}

// SendResult submits the result of a computed task
func (c *Client) SendResult(id string, result float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.cli.SubmitTask(ctx, &pb.TaskResult{Id: id, Result: result})
	return err
}
