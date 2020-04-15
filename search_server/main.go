//go:generate protoc -I ../proc --go_out=grpc:../proc ../proc/search.proto
package main

import (
	"context"
	"google.golang.org/grpc"
	"hello/proc"
	"log"
	"net"
	"strconv"
)

type server struct {
	proc.UnimplementedSearchServiceServer
}
const (
	port = ":50051"
)

func (s *server) Search(ctx context.Context, in *proc.SearchRequest) (*proc.SearchResponse, error) {
	log.Printf("Received: %v", in.GetPageNumber())
	return &proc.SearchResponse{Message: strconv.Itoa(int(in.GetPageNumber()))}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proc.RegisterSearchServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
