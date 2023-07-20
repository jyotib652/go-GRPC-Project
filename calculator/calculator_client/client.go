package main

import (
	"context"
	"fmt"
	"log"
	"myapp/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	fmt.Println("I'm a client.")
	// conn, err := grpc.Dial("the address to connect to", bunch_of_options)
	// cc for client connection
	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Couldn't connect to server: ", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	doUnary(c)

	// fmt.Printf("Created the client %f\n", c)
	// req := &calculatorpb.SumRequest{
	// 	FirstNumber:  3,
	// 	SecondNumber: 10,
	// }
	// resp, err := c.Sum(context.Background(), req)
	// if err != nil {
	// 	log.Fatal("error while calling Greet RPC:", err)
	// }
	// log.Println("Response from Greet:", resp.SumResult)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sum unary RPC...")
	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 10,
	}
	resp, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling Greet RPC:", err)
	}
	log.Println("Response from Greet:", resp.SumResult)
}
