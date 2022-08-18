package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	apiv1 "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/fbngrm/bid-tracker/pkg/item"
	"github.com/fbngrm/bid-tracker/pkg/user"
	"github.com/fbngrm/bid-tracker/server/internal/api"
	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
}

// NewServer returns an Server instance with a service attached.
func NewServer(ctx context.Context) (*Server, error) {
	// deps
	itemService := item.NewService()
	userService := user.NewService()
	bidService := bid.NewService(
		itemService,
		userService,
	)
	handler := api.NewApi(bidService, itemService, userService)

	// note, not a production ready config
	server := grpc.NewServer()
	apiv1.RegisterServiceServer(server, handler)

	// seed
	userService.

	return &Server{
		server: server,
	}, nil
}

func (s *Server) Run(grpcEndpoint string) error {
	// TODO: add config
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", grpcEndpoint, err)
	}

	log.Println("gRPC listening on", grpcEndpoint)
	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}
	return nil
}

// stops the gRPC server gracefully. It stops the server from accepting new
// connections and RPCs and blocks until all the pending RPCs are  finished.
func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
