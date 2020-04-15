package main

import (
	"context"
	"google.golang.org/grpc"
	"hello/proc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proc.NewSearchServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Search(ctx, &proc.SearchRequest{PageNumber: 100})
	if err != nil {
		log.Fatalf("could not search: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
