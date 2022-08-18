package server

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fbngrm/bid-tracker/server/internal/grpc"
	"github.com/fbngrm/bid-tracker/server/internal/http"
)

var (
	grpcEndpoint string
)

func main() {
	flag.StringVar(&grpcEndpoint, "grpc-endpoint", ":8080", "gRPC server endpoint")
	flag.Parse()

	// todo, add redis as event store for aggregates

	ctx, cancel := context.WithCancel(context.Background())

	quitCh := make(chan os.Signal, 1)
	// interrupt signal sent from terminal
	signal.Notify(quitCh, os.Interrupt)
	// sigterm signal sent from kubernetes or docker
	signal.Notify(quitCh, syscall.SIGTERM)

	go func() {
		sig := <-quitCh
		log.Printf("received signal %s:, shutting down...\n", sig)
		cancel()
	}()

	// http
	httpServer, err := http.NewServer(ctx, grpcEndpoint)
	if err != nil {
		log.Printf("could not init http server: %v\n", err)
		os.Exit(1)
	}
	go func() {
		if err := httpServer.Run(); err != nil {
			log.Printf("error running http server: %v\n", err)
		}
	}()

	// grpc
	grpcServer, err := grpc.NewServer(ctx)
	if err != nil {
		log.Printf("could not init grpc server: %v\n", err)
		os.Exit(1)
	}
	go func() {
		if err := grpcServer.Run(grpcEndpoint); err != nil {
			log.Printf("error running grpc server: %v\n", err)
		}
	}()

	<-ctx.Done() // wait for shutdown signal

	// we want to give the shutdown a timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)

	go func() {
		<-shutdownCtx.Done()
		if errors.Is(shutdownCtx.Err(), context.Canceled) {
			return
		}
		log.Printf("couldn't shutdown gracefully: %v", shutdownCtx.Err())
		os.Exit(1)
	}()

	// shutdown servers
	graceful := true
	if err := httpServer.Shutdown(ctx); err != nil {
		graceful = false
		log.Printf("error shutting down http server: %q\n", err)
	}
	grpcServer.Shutdown()

	shutdownCancel()
	log.Println("shutdown complete")

	if graceful {
		os.Exit(0)
	}
	os.Exit(1)
}
