package main

import (
	"context"
	"github.com/stashchenko/microservice-example/pkg/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := proto.NewHealthClient(conn)
	cl, err := client.Watch(ctx, &proto.HealthCheckRequest{
		Service: "test",
	})
	for {
		resp, err := cl.Recv()
		log.Println("Res:", resp, "Error:", err)
		time.Sleep(1 * time.Second)
	}
}
