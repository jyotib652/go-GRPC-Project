package main

import (
	"context"
	"fmt"
	"log"
	"myapp/greet/greetpb"

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

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)

	// fmt.Printf("Created the client %f\n", c)
	// req := &greetpb.GreetRequest{
	// 	Greeting: &greetpb.Greeting{
	// 		FirstName: "Jyotiprokash",
	// 		LastName:  "Ban",
	// 	},
	// }
	// resp, err := c.Greet(context.Background(), req)
	// if err != nil {
	// 	log.Fatal("error while calling Greet RPC:", err)
	// }
	// log.Println("Response from Greet:", resp.Result)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Greet unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jyotiprokash",
			LastName:  "Ban",
		},
	}
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatal("error while calling Greet RPC:", err)
	}
	log.Println("Response from Greet:", resp.Result)
}
