package main

import (
	"context"
	"fmt"
	"log"
	"myapp/greet/greetpb"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Greet function was invoked with:", req)
	first_anme := req.GetGreeting().GetFirstName()
	result := "Hello, " + first_anme
	resp := &greetpb.GreetResponse{
		Result: result,
	}
	return resp, nil
}

func main() {
	fmt.Println("Greet Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	fmt.Println("listenning on port 50051")

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
