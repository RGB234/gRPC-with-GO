package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/rgb234/gRPC-with-GO/clientstreaming/pbs"
	"google.golang.org/grpc"
)

var (
	// flag name, default value, usage message (--help)
	port = flag.Int("port", 50051, "The server port")
)

type server struct{
	pb.UnimplementedClientStreamingServer
}

func (s *server) ProcessIO(stream pb.ClientStreaming_ProcessIOServer) error {
	fmt.Printf("Server processing gRPC client-streaming.\n")
	// 종료 조건까지 무한루프
	var count int32 = 0;
	for {
		// Receive messages from client
		msg, err := stream.Recv() 
		if err == io.EOF { 
			// no more messages
			count += 1
			return stream.SendAndClose(&pb.Number{Value: count})
		  }
		if err != nil {
			log.Fatalf("failed to Recv: %v", err)
			return err
		}
		log.Printf("[received from client] %s", msg)
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		// 에러 , 로그 출력 후 종료
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterClientStreamingServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}