package api

import (
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	userpb "example/internal/rpc/user/pb"
	serviceUser "example/internal/services/user"
)

// New builds a gRPC server with every configured resource's service
// registered, wired to its domain repository and database connection.
func New(db *sqlx.DB) *grpc.Server {
	s := grpc.NewServer()

	repo := serviceUser.NewRepository(db)
	svc := serviceUser.NewService(repo)
	userpb.RegisterUserServiceServer(s, svc)

	return s
}
