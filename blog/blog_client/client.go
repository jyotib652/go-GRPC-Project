package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"myapp/blog/blogpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	fmt.Println("Blog Client")
	// conn, err := grpc.Dial("the address to connect to", bunch_of_options)
	// cc for client connection

	// ssl implementation for cliet
	opts := []grpc.DialOption{}
	certfile := "ssl/ca.crt" // Certificate Authority Trust Certificate
	creds, sslErr := credentials.NewClientTLSFromFile(certfile, "")
	if sslErr != nil {
		log.Fatal("Error while loading CA Trust Certificate: ", sslErr)
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))
	cc, err := grpc.Dial("localhost:50051", opts...)
	// cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Couldn't connect to server: ", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Creating the Blog
	blog := &blogpb.Blog{
		AuthorId: "Jyotiprokash",
		Title:    "My First Blog",
		Content:  "content of the first blog.",
	}

	createBlogResp, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatal("Unexpected error while calling CreateBlog() from client: ", err)
	}

	fmt.Println("Blog has been created: ", createBlogResp)
	blogID := createBlogResp.GetBlog().GetId()

	// Read the blog
	fmt.Println("Reading the blog")
	// Here, we are trying to read random blog which doesn't exist
	// Because {BlogId: "138947957975"} doesn't exist in the database so, it's going to fail
	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "138947957975"})
	if err != nil {
		fmt.Println("Error occured while reading from the database", err)
	}

	// This time blogID really exist in the mongodb database so, it's going to work
	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogResp, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		fmt.Println("Error occured while reading from the database", err)
	}

	fmt.Println("A Blog was read: ", readBlogResp)

	// Update a Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "John Smith",
		Title:    "My First Blog (edited)",
		Content:  "content of the first blog, with some awesome additions!",
	}

	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Println("Error occured while updating the database :", err)
	}

	fmt.Println("Blog was updated :", updateRes)

	// delete Blog
	// Here, we are passing previously created `blogID`
	deleteResp, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Println("Error occured while deleting from the database :", err)
	}

	fmt.Println("Blog was deleted : ", deleteResp)

	// list Blogs

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if err != nil {
		log.Fatal("Error while calling ListBlog RPC: ", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error occured: %v \n", err)
		}
		fmt.Println(res.GetBlog())
	}

}
