package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/rgb234/gRPC-with-GO/serverstreaming/pbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr =  flag.String("addr", "localhost:50051", "the address to connect to")
)

func recv_message(client pb.ServerStreamingClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// send to server
	request := pb.Number{Value: 5}
	stream, err := client.ProcessIO(ctx, &request)
	if err != nil {
		log.Fatalf("[failed to creating stream] : %v", err)
	}
	// Receive from server
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			// receive end
			break;
		}
		if err != nil {
			log.Fatalf("[Recv error] : %v", err)
		}
		log.Printf("[server to client] %s", response.GetMessage())
	}
	
}

func main() {
	conn,err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewServerStreamingClient(conn)
	recv_message(client)
}