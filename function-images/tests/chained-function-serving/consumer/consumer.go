package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "tests/chained-functions-serving/proto"
)

type consumerServer struct {
	pb.UnimplementedProducerConsumerServer
}

func (s *consumerServer) ConsumeString(ctx context.Context, str *pb.ConsumeStringRequest) (*pb.ConsumeStringReply, error) {
	log.Printf("[consumer] Consumed %v\n", str.Value)
	return &pb.ConsumeStringReply{Value: true}, nil
}

func (s *consumerServer) ConsumeStream(stream pb.ProducerConsumer_ConsumeStreamServer) error {
	for {
		str, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ConsumeStringReply{Value: true})
		}
		if err != nil {
			return err
		}
		log.Printf("[consumer] Consumed %v\n", str.Value)
	}
}

func main() {
	port := flag.Int("p", 3030, "Port")
	flag.Parse()

	//set up log file
	file, err := os.OpenFile("log/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	//set up server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("[consumer] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := consumerServer{}
	pb.RegisterProducerConsumerServer(grpcServer, &s)

	log.Println("[consumer] Server Started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[consumer] failed to serve: %s", err)
	}

}
