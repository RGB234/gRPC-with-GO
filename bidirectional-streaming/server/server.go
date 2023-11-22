package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/rgb234/gRPC-with-GO/bidirectional-streaming/pbs"
	"google.golang.org/grpc"
)

type server struct{
	pb.UnimplementedBidirectionalServer
}

func (s *server) GetServerResponse(stream pb.Bidirectional_GetServerResponseServer) error {
	fmt.Printf("Server processing gRPC bidirectional streaming.")
	for {
		message, err := stream.Recv()
		if err == nil {
			return err
		}
		fmt.Println(message)
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		// 에러 , 로그 출력 후 종료
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBidirectionalServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}