package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/rgb234/gRPC-with-GO/serverstreaming/pbs"
	"google.golang.org/grpc"
)

var (
	// flag name, default value, usage message (--help)
	port = flag.Int("port", 50051, "The server port")
)

type server struct{
	pb.UnimplementedServerStreamingServer
}

func make_message(message string) *pb.Message {
	return &pb.Message{Message: message}
}

func (s *server) ProcessIO(number *pb.Number, stream pb.ServerStreaming_ProcessIOServer) error {
	// receive from client
	num := number.GetValue()
	fmt.Printf("Server processing gRPC bidirectional streaming. {%d}\n", num)
	
	// send messages to client
	messages := []*pb.Message  {
		make_message("message #1"),
		make_message("message #2"),
		make_message("message #3"),
		make_message("message #4"),
		make_message("message #5"),
	}
	for _, msg :=  range(messages) {
		err := stream.Send(&pb.Message{Message: msg.GetMessage()})
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("failed to Send: %v", err)
			return err
		}
		log.Printf("[send to client] %s", msg)
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		// 에러 , 로그 출력 후 종료
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterServerStreamingServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}