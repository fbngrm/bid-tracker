package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	gw "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO: add config
const port = "8081"
const openApiPrefix = "/openapi"

type Server struct {
	server *http.Server
}

// NewServer returns an HTTPServer instance with a handler attached.
// Note, handlers should implement a timeout to avoid running into transport
// layer timeouts.
// Add HTTP middleware here.
func NewServer(ctx context.Context, grpcEndpoint string) (*Server, error) {
	mux := runtime.NewServeMux()

	// openapi documentation handler
	fs := http.FileServer(http.Dir("./apis/ui"))
	if err := mux.HandlePath("GET", openApiPrefix+".json", handleOpenapiDescription); err != nil {
		return nil, fmt.Errorf("could not register openapi json endpoint: %v", err)
	}
	if err := mux.HandlePath("GET", openApiPrefix+"/*", handler(http.StripPrefix(openApiPrefix, fs))); err != nil {
		return nil, fmt.Errorf("could not register openapi endpoint: %v", err)
	}

	// regiester with gRPC endpoint
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterMatchServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("could not register mux router with gRPC endpoint: %v", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second, // deadline for reading request body
		WriteTimeout: 30 * time.Second, // deadline for ServeHTTP
	}
	return &Server{server: server}, nil

}

func (s *Server) Run() error {
	log.Printf("HTTP listening on :%s\n", port)
	return s.server.ListenAndServe()
}

// Shutdown stops accepting new requests and waits for the running
// ones to finish before returning. See net/http docs for details.
// The provided context should have a timeout.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// implements mux.HandlerFunc, we ignore path params in this scenario.
func handler(h http.Handler) func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	return func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		h.ServeHTTP(w, r)
	}
}

// implements mux.HandlerFunc, we ignore path params in this scenario.
func handleOpenapiDescription(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	content, err := os.ReadFile("gen/openapiv2/auction/v1/auction.swagger.json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	if _, err = w.Write(content); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
