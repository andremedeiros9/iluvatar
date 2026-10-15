package grpcserver

import (
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	"manwe/internal/user"
	userpb "manwe/internal/user/pb"
)

// New builds a gRPC server with every configured resource's service
// registered, backed by db.
func New(db *sqlx.DB) *grpc.Server {
	s := grpc.NewServer()

	userpb.RegisterUserServiceServer(s, user.NewService(user.NewRepository(db)))

	return s
}
