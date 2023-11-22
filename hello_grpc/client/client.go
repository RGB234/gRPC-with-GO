package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/rgb234/gRPC-with-GO/hello_grpc/pbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inputVal int32 = 4;
)

var (
	addr =  flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	conn,err := grpc.Dial(*addr,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil{
		log.Fatalf("connection failed: %v", err)
	}
	// main 함수 종료 직전에 연결 해제
	defer conn.Close()
	// client stub
	c := pb.NewMyServiceClient(conn)

	// Contact the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Remote procedure call
	result, err := c.MyFunction(ctx, &pb.MyNumber{Value: inputVal})
	if err != nil{
		log.Fatalf("rpc failed: %v", err)
	}
	log.Printf("%d^2 = %d", inputVal, result.GetValue())
}