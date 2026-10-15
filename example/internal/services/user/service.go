package user

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	domainuser "example/internal/domain/user"
	pb "example/internal/rpc/user/pb"
)

// Service implements pb.UserServiceServer.
type Service struct {
	pb.UnimplementedUserServiceServer

	repo domainuser.Repository
}

// NewService returns a Service backed by repo.
func NewService(repo domainuser.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	m := &domainuser.User{}

	m.Email = req.Email

	m.Age = req.Age

	if err := s.repo.Create(ctx, m); err != nil {
		return nil, status.Errorf(codes.Internal, "creating user: %v", err)
	}

	return &pb.CreateUserResponse{User: toProto(m)}, nil
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	m, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user %s not found: %v", req.Id, err)
	}

	return &pb.GetUserResponse{User: toProto(m)}, nil
}

func (s *Service) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "listing users: %v", err)
	}

	resp := &pb.ListUserResponse{}
	for i := range items {
		resp.Users = append(resp.Users, toProto(&items[i]))
	}
	return resp, nil
}

func (s *Service) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	m := &domainuser.User{ID: id}

	m.Email = req.Email

	m.Age = req.Age

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, status.Errorf(codes.Internal, "updating user %s: %v", req.Id, err)
	}

	return &pb.UpdateUserResponse{User: toProto(m)}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id: %v", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "deleting user %s: %v", id, err)
	}

	return &pb.DeleteUserResponse{}, nil
}

func toProto(m *domainuser.User) *pb.User {
	return &pb.User{
		Id:        m.ID.String(),
		Email:     m.Email,
		Age:       m.Age,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}
}
