package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"

	pb "tests/chained-functions-serving/proto"
)

func main() {
	strings := flag.Int("s", 1, "Number of strings to send")
	address := flag.String("addr", "localhost", "Server IP address")
	clientPort := flag.Int("pc", 3031, "Client Port")
	flag.Parse()

	fmt.Printf("Client using address: %v\n", *address)

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *address, *clientPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %s", err)
	}
	defer conn.Close()

	client := pb.NewClientProducerClient(conn)
	client.ProduceStrings(context.Background(), &pb.ProduceStringsRequest{Value: int32(*strings)})

	fmt.Printf("client closing")

}
