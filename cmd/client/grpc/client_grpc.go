package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/anras5/todo-app-backend/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	requestCount = 100000
	grpcURL      = "localhost:9000"
)

func main() {

	conn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTodoServiceClient(conn)

	// Create
	startTime := time.Now()
	ids := []int{}
	for i := 0; i < requestCount; i++ {
		createdTodo, err := client.Create(context.Background(), &pb.Todo{
			Name:        fmt.Sprintf("Todo number %d", i),
			Description: "This is a todo",
			Deadline:    timestamppb.New(time.Now().Add(24 * time.Hour)),
			Completed:   false,
		})
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", i, err)
		}
		id := int(createdTodo.Id)
		ids = append(ids, id)
	}
	duration := time.Since(startTime)
	fmt.Printf("Create %v.\n", duration)

	// Get by one
	startTime = time.Now()
	for _, id := range ids {
		_, err := client.Get(context.Background(), &pb.Id{Id: int32(id)})
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("Get %v.\n", duration)

	// Update
	startTime = time.Now()
	for i, id := range ids {
		_, err := client.Update(context.Background(), &pb.Todo{
			Id:          int32(id),
			Name:        fmt.Sprintf("Todo number %d", i),
			Description: "This is an updated todo",
			Deadline:    timestamppb.New(time.Now().Add(24 * time.Hour)),
			Completed:   false,
		})
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("Update %v.\n", duration)

	// List
	startTime = time.Now()
	stream, err := client.List(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Printf("Error sending request for Todo list: %v\n", err)
	}
	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("List %v.\n", duration)

	// Delete
	startTime = time.Now()
	for _, id := range ids {
		_, err := client.Delete(context.Background(), &pb.Id{Id: int32(id)})
		if err != nil {
			fmt.Printf("Error sending request for Todo #%d: %v\n", id, err)
		}
	}
	duration = time.Since(startTime)
	fmt.Printf("Delete %v.\n", duration)

}
