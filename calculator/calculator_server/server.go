package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"myapp/calculator/calculatorpb"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Recieved Sum RPC:", req)
	first_number := req.FirstNumber
	second_number := req.SecondNumber
	sum := first_number + second_number
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}

	return res, nil
}

func main() {
	// Before => There is no implementation of logs to get the filename and file number upon crashing go code
	// After => Setting the log level to something different for Errors
	// If we crash the go code, we get the file name and line nuber
	// Begin
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// End

	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	fmt.Println("listenning on port 50051")

	// Before => Simple implementation without closing of server and listener
	// Begin
	// s := grpc.NewServer()
	// calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	//
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatal("failed to serve:", err)
	// }
	//
	// End

	// After => Better implementation with closing of server & listener upon "Ctrl+ c"
	// Begin
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// STOPPING THE SERVER GRACEFULLY
	go func() {
		fmt.Println("Starting the server")
		if err := s.Serve(lis); err != nil {
			log.Fatal("failed to serve:", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("End of Program")
	// End
}
