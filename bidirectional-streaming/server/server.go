package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"time"

	pb "github.com/rgb234/gRPC-with-GO/bidirectional-streaming/pbs"
	"google.golang.org/grpc"
)

var (
	// flag name, default value, usage message (--help)
	port = flag.Int("port", 50051, "The server port")
)

type server struct{
	pb.UnimplementedBidirectionalServer
}

func (s *server) ProcessIoStream(stream pb.Bidirectional_ProcessIoStreamServer) error {
	fmt.Printf("Server processing gRPC bidirectional streaming.\n")
	// 종료 조건까지 무한루프
	for {
		// Receive messages from client
		msg, err := stream.Recv() 
		if err == io.EOF { 
			// no more messages
			return nil
		  }
		if err != nil {
			log.Fatalf("failed to Recv: %v", err)
			return err
		}
		log.Printf("[received from client] %s", msg)

		// processing delay
		time.Sleep(1 * time.Second)

		// send messages to client
		err = stream.Send(&pb.Message{Message: msg.GetMessage()})
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("failed to Send: %v", err)
			return err
		}
		log.Printf("[send to client] %s", msg)
	}
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Printf("current cpu max num : %d \n", numCPU)
	
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
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