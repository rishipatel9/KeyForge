package main

import (
	"context"
	"fmt"
	"kvstorepb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	kvstorepb.UnimplementedKeyValueStoreServer
	store map[string]string
}

func (s *server) Set(ctx context.Context, kv *kvstorepb.KeyValue) (*kvstorepb.Response, error) {
	s.store[kv.Key] = kv.Value
	return &kvstorepb.Response{Message: "Key-Value set successfully!"}, nil
}

func (s *server) Get(ctx context.Context, key *kvstorepb.Key) (*kvstorepb.Response, error) {
	value, exists := s.store[key.Key]
	if !exists {
		return &kvstorepb.Response{Message: "Key not found"}, nil
	}
	return &kvstorepb.Response{Message: "Success", Value: value}, nil
}

func startServer() {
	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	kvstorepb.RegisterKeyValueStoreServer(grpcServer, &server{store: make(map[string]string)})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func main() {
	fmt.Println("Starting gRPC server...")
	startServer()
}
