/*
ref. https://grpc.io/docs/languages/go/basics/
ref. https://github.com/grpc/grpc-go//examples/route_guide

Note that in gRPC-Go, RPCs operate in a blocking/synchronous mode,
which means that the RPC call waits for the server to respond

Although each side will always get the other’s messages in the order they were written,
both the client and server can read and write in any order — the streams operate completely independently.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/rgb234/gRPC-with-GO/bidirectional-streaming/pbs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr =  flag.String("addr", "localhost:50051", "the address to connect to")
)


func make_message(message string) *pb.Message {
	return &pb.Message{Message:message}
}

func generate_messages() chan *pb.Message {
	messages := []string{
		"message #1",
		"message #2",
		"message #3",
		"message #4",
		"message #5",
	}
	// buffered channel, buffer size : len(messages)
	messageChan := make(chan *pb.Message, len(messages))

	// yield
	// ref. http://golang.site/go/article/22-Go-%EC%B1%84%EB%84%90
	// 익명함수 정의, 별도의 goroutine 으로 익명함수 실행
	go func() { 
		// 해당 익명함수 종료시 (for 구문으로 모든 메시지를 전송했을 시) 채널닫기
		defer close(messageChan) 
		// channel range 구문
		for _, message := range messages { 
			messageChan <- make_message(message)
		}
	// () : 해당 익명함수에 파라미터 전달
	// ref. http://golang.site/go/article/21-Go-%EB%A3%A8%ED%8B%B4-goroutine
	}() 
	return messageChan
}

func message_handler(client pb.BidirectionalClient, messageChan chan *pb.Message){
	// 서버로부터 메시지 수신 & 서버로 메시지 전송
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	waitc := make(chan struct{})
	stream, err := client.GetServerResponse(ctx)
	if err != nil {
		log.Fatalf("failed to creating a stream: %v", err)
	}
	//server to client
	go func(){
		 // 종료조건까지 무한 반복
		for {
			 // Receive from server
			msg, err := stream.Recv()
			if err == io.EOF {
				// read done
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive a message from the server")
			}
			fmt.Printf("[received from server] %s\n", msg)
		}
	}()
	// send to server
	for msg := range messageChan {
		fmt.Printf("[send to server] %s \n", msg)
		err := stream.Send(&pb.Message{Message: msg.GetMessage()}) // send messages to server
		if err == io.EOF {
			// send done
			return
		}
		if err != nil{
			log.Fatalf("failed to send a message to the server")
		}
	}
	stream.CloseSend()
	// waitc 채널이 닫히지 않는 이상 waitc 채널의 값을 수신하기 위해 계속대기 (blocking)
	// 수신이 모두 완료될 경우 (err == io.EOF) 채널이 닫히면서 대기종료 
	<-waitc
}

func main(){
	conn,err := grpc.Dial(*addr,grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBidirectionalClient(conn)
	message_handler(client, generate_messages())
}