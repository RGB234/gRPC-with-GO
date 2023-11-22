package main

import (
	"fmt"

	pb "github.com/rgb234/gRPC-with-GO/bidirectional-streaming/pbs"
)

func make_message(message string) pb.Message {
	return pb.Message{Message:message}
}

func generate_messages() pb.Message {
	var messages []pb.Message = []pb.Message {
		make_message("message #1"),
		make_message("message #2"),
		make_message("message #3"),
		make_message("message #4"),
		make_message("message #5"),
	}
	for _, msg := range(messages) {
		fmt.Printf("[client to server] %v", msg.Message)
		return msg.Message
	}
}

func send_message()

func main(){

}