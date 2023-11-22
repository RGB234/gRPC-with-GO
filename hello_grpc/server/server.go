package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	hello_grpc "github.com/rgb234/gRPC-with-GO/hello_grpc"
	pb "github.com/rgb234/gRPC-with-GO/hello_grpc/pbs"
	"google.golang.org/grpc"
)

var (
	// flag name, default value, usage message (--help)
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedMyServiceServer
}

func (s *server) MyFunction(ctx context.Context, in *pb.MyNumber) (*pb.MyNumber, error){
	log.Printf("Received: %d", in.GetValue())
	result := hello_grpc.My_func(in.GetValue())
	return &pb.MyNumber{Value: result}, nil
}

func main(){
	// cmd 명령어 파싱
	flag.Parse()
	// Specify the port we want to use to listen for client requests using
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		// 로그 출력 후 종료
		log.Fatalf("failed to listen: %v", err)
	}
	// new Server Instance 생성 & 추가
	// concurrency : ref. https://github.com/grpc/grpc-go/blob/master/Documentation/concurrency.md
	// Each RPC handler attached to a registered server will be invoked in its own goroutine.
	// multiple services can be registered to the same server.
	
	// Gorountines : ref. https://go.dev/tour/concurrency/1 
	// lightweight thread managed by the Go runtime
	// server (struct) 에 등록된 각 grpcServer 객체마다 각각 gorountine 을 가짐 -> 알아서 멀티스레딩
	grpcServer := grpc.NewServer()
	pb.RegisterMyServiceServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	// Serve() 실행, Serve() 실패시 log 출력
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}