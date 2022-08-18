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
	u1, err := userService.CreateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not seed user: %w", err)
	}
	u2, err := userService.CreateUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not seed user: %w", err)
	}
	i1, err := itemService.CreateItem(ctx, "picasso")
	if err != nil {
		return nil, fmt.Errorf("could not seed items: %w", err)
	}
	i2, err := itemService.CreateItem(ctx, "beltracci")
	if err != nil {
		return nil, fmt.Errorf("could not seed items: %w", err)
	}
	fmt.Println("======================================================")
	fmt.Println("SEED DATA - USE FOR INTERACTING WITH THE API INITIALLY")
	fmt.Println("------------------------------------------------------")
	fmt.Printf("user 1: id = %s\n", u1.ID.String())
	fmt.Printf("user 2: id = %s\n", u2.ID.String())
	fmt.Println("------------------------------------------------------")
	fmt.Printf("item 1: id = %s\n", i1.ID.String())
	fmt.Printf("item 2: id = %s\n", i2.ID.String())
	fmt.Println("------------------------------------------------------")
	fmt.Println("Place a bid:")
	fmt.Println(fmt.Sprintf(`curl -d '{"item_id":"%s", "user_id":"%s", "amount":1.5}' http://localhost:8081/v1/bid`, i1.ID.String(), u1.ID.String()))
	fmt.Println("Get highest bid for item:")
	fmt.Printf("curl http://localhost:8081/v1/items/%s/bids/highest\n", i1.ID.String())
	fmt.Println("Get bids for item:")
	fmt.Printf("curl http://localhost:8081/v1/item/%s/bids\n", i1.ID.String())
	fmt.Println("Get items for users bids:")
	fmt.Printf("curl http://localhost:8081/v1/user/%s/bids/items\n", u1.ID.String())
	fmt.Println("======================================================")

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
