package rpc

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"

	"github.com/anras5/todo-app-backend/internal/grpc/pb"
	"github.com/anras5/todo-app-backend/internal/models"
	"github.com/anras5/todo-app-backend/internal/repository"
	"github.com/anras5/todo-app-backend/internal/repository/dbrepo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TodoServer struct {
	DB repository.DatabaseRepo
}

func NewTodoServer(db *sql.DB) *TodoServer {
	return &TodoServer{
		DB: dbrepo.NewPostgresRepo(db),
	}
}

func (s *TodoServer) Run() {
	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTodoServiceServer(grpcServer, s)

	log.Println("starting grpc server on :9000")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc server over 9000: %v", err)
	}
}

func (s *TodoServer) Create(ctx context.Context, req *pb.Todo) (*pb.Todo, error) {
	todo := models.Todo{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Deadline:    req.GetDeadline().AsTime(),
		Completed:   req.GetCompleted(),
	}

	id, err := s.DB.InsertTodo(todo)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.Todo{
		Id:          int32(id),
		Name:        todo.Name,
		Description: todo.Description,
		Deadline:    timestamppb.New(todo.Deadline),
		Completed:   todo.Completed,
	}, nil
}

func (s *TodoServer) Get(ctx context.Context, req *pb.Id) (*pb.Todo, error) {
	id := int(req.GetId())

	todo, err := s.DB.SelectTodo(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "todo not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.Todo{
		Id:          int32(todo.ID),
		Name:        todo.Name,
		Description: todo.Description,
		Deadline:    timestamppb.New(todo.Deadline),
		Completed:   todo.Completed,
	}, nil
}

func (s *TodoServer) Update(ctx context.Context, req *pb.Todo) (*pb.Todo, error) {
	todo := models.Todo{
		ID:          int(req.GetId()),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Deadline:    req.GetDeadline().AsTime(),
		Completed:   req.GetCompleted(),
	}

	err := s.DB.UpdateTodo(todo)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.Todo{
		Id:          int32(todo.ID),
		Name:        todo.Name,
		Description: todo.Description,
		Deadline:    timestamppb.New(todo.Deadline),
		Completed:   todo.Completed,
	}, nil
}

func (s *TodoServer) Delete(ctx context.Context, req *pb.Id) (*pb.Todo, error) {
	id := int(req.GetId())

	todo, err := s.DB.SelectTodo(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "todo not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	err = s.DB.DeleteTodo(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "todo not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.Todo{
		Id:          int32(todo.ID),
		Name:        todo.Name,
		Description: todo.Description,
		Deadline:    timestamppb.New(todo.Deadline),
		Completed:   todo.Completed,
	}, nil
}

func (s *TodoServer) List(t *emptypb.Empty, stream grpc.ServerStreamingServer[pb.Todo]) error {
	todos, err := s.DB.SelectTodos()
	if err != nil {
		return status.Error(codes.Internal, "internal error")
	}

	for _, todo := range todos {
		err := stream.Send(&pb.Todo{
			Id:          int32(todo.ID),
			Name:        todo.Name,
			Description: todo.Description,
			Deadline:    timestamppb.New(todo.Deadline),
			Completed:   todo.Completed,
		})
		if err != nil {
			return status.Error(codes.Internal, "internal error")
		}
	}

	return nil
}
