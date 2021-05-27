package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"

	pb "tests/chained-functions-serving/proto"
)

type producerServer struct {
	producerClient pb.ProducerConsumerClient
	pb.UnimplementedClientProducerServer
}

func (ps *producerServer) ProduceStrings(c context.Context, count *pb.ProduceStringsRequest) (*pb.Empty, error) {
	if count.Value <= 0 {
		return new(pb.Empty), nil
	} else if count.Value == 1 {
		produceSingleString(ps.producerClient, fmt.Sprint(rand.Intn(1000)))
	} else {
		wordList := make([]string, int(count.Value))
		for i := 0; i < int(count.Value); i++ {
			wordList[i] = fmt.Sprint(rand.Intn(1000))
		}
		produceStreamStrings(ps.producerClient, wordList)
	}
	return new(pb.Empty), nil
}

func produceSingleString(client pb.ProducerConsumerClient, s string) {
	ack, err := client.ConsumeString(context.Background(), &pb.ConsumeStringRequest{Value: s})
	if err != nil {
		log.Fatalf("client error in string consumption: %s", err)
	}
	fmt.Printf("(single) Ack: %v\n", ack.Value)
}

func produceStreamStrings(client pb.ProducerConsumerClient, strings []string) {
	//make stream
	stream, err := client.ConsumeStream(context.Background())
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
	}

	//stream strings
	for _, s := range strings {
		if err := stream.Send(&pb.ConsumeStringRequest{Value: s}); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, s, err)
		}
	}

	//end transaction
	ack, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	fmt.Printf("(stream) Ack: %v\n", ack.Value)
}

func main() {
	address := flag.String("addr", "localhost", "Server IP address")
	clientPort := flag.Int("pc", 3030, "Client Port")
	serverPort := flag.Int("ps", 3031, "Server Port")
	flag.Parse()

	//client setup
	fmt.Printf("Client using address: %v\n", *address)

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", *address, *clientPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %s", err)
	}
	defer conn.Close()

	client := pb.NewProducerConsumerClient(conn)

	//produceSingleString(client, "hello")
	//strings := []string{"Hello", "World", "one", "two", "three"}
	//produceStreamStrings(client, strings)

	//server setup
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *serverPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := producerServer{}
	s.producerClient = client
	pb.RegisterClientProducerServer(grpcServer, &s)

	fmt.Println("Server Started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
