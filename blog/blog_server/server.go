package main

import (
	"context"
	"fmt"
	"log"
	"myapp/blog/blogpb"
	"os"
	"os/signal"
	"time"

	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct {
	// blogpb.UnimplementedBlogServiceServer
	blogpb.BlogServiceServer
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Printf("received a request to Create a Blog with parameters: %v\n", req)
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	// Now, we'll insert this data to our mongodb database
	res, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal database error: %v", err),
		)
	}
	// Now, we are extracting the object id of the inserted data; oid for objectid
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			"Unable to extract id of inserted data object or couldn't convert it to oid",
		)
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Printf("received a request to Read a Blog with parameters: %v\n", req)
	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Couldn't parse ID: %v", err),
		)
	}

	// create an empty struct
	data := &blogItem{}
	// filter := bson.NewDocument(bson.EC.ObjectID("_id", oid))
	filter := bson.M{"_id": oid}
	res := collection.FindOne(context.Background(), filter)
	err = res.Decode(&data)
	if err != nil {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("Could not find blog with specified ID,:  %v", err),
		)
	}

	return &blogpb.ReadBlogResponse{
		// Blog: &blogpb.Blog{
		// 	Id:       data.ID.Hex(),
		// 	AuthorId: data.AuthorID,
		// 	Title:    data.Title,
		// 	Content:  data.Content,
		// },

		Blog: dataToBlogPb(data),
	}, nil
}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}
}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Printf("received a request to Update a Blog with parameters: %v\n", req)
	blog := req.GetBlog()
	// get object id
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Couldn't parse ID: %v\n", err),
		)
	}

	// // Method-1 : first read, then update the fields and then replace the actual value (old method)
	// // BEGIN
	// // create an empty struct
	// data := &blogItem{}
	// // filter := bson.NewDocument(bson.EC.ObjectID("_id", oid))
	// filter := bson.M{"_id": oid}
	// res := collection.FindOne(context.Background(), filter)
	// err = res.Decode(&data)
	// if err != nil {
	// 	return nil, status.Error(
	// 		codes.NotFound,
	// 		fmt.Sprintf("Could not find blog with specified ID,:  %v", err),
	// 	)
	// }

	// // we updated our internal struct
	// data.AuthorID = blog.GetAuthorId()
	// data.Title = blog.GetTitle()
	// data.Content = blog.GetContent()

	// // updateRes, err := collection.ReplaceOne(ctx, filter, data)
	// _, err = collection.ReplaceOne(ctx, filter, data)
	// if err != nil {
	// 	return nil, status.Error(
	// 		codes.Internal,
	// 		fmt.Sprintf("Could not update the object in mongodb, : %v", err),
	// 	)
	// }
	// // END

	// Method-2
	// BEGIN
	data := &blogItem{
		AuthorID: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}

	filter := bson.M{"_id": oid}
	res, err := collection.UpdateOne(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Could not update: %v\n", err),
		)
	}

	if res.MatchedCount == 0 {
		// this means we didn't find any entry with the provided id
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find the blog with the Id: %v\n", err),
		)
	}
	// END

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Printf("received a request to Delete a Blog with parameters: %v\n", req)
	// get object id
	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Couldn't parse ID: %v", err),
		)
	}

	filter := bson.M{"_id": oid}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("Could not delete the object in mongodb, : %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Error(
			codes.NotFound,
			fmt.Sprintf("Could not find the Blog in mongodb, : %v", err),
		)
	}

	return &blogpb.DeleteBlogResponse{BlogId: req.GetBlogId()}, nil
}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("received a request to List a Blog")

	// Since, we need to get all the Blogs in our database, we don't need to
	// provide a filter. So we are providing nil for filter. cur for cursor
	// cur, err := collection.Find(context.Background(), nil)
	cur, err := collection.Find(context.Background(), bson.D{})

	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error : %v", err),
		)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &blogItem{}
		err = cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error occured while decoding cursor data :%v", err),
			)
		}
		stream.Send(&blogpb.ListBlogResponse{Blog: dataToBlogPb(data)})
	}

	err = cur.Err()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal cursor error: %v", err),
		)
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to Mongodb")

	// Connecting to MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:secret321@localhost:27017"))
	if err != nil {
		log.Fatal("Couldn't establish a connection to MongoDB server:", err)
	}

	// Make sure to defer a call to Disconnect after instantiating your client:
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Calling Connect does not block for server discovery.
	// If you wish to know if a MongoDB server has been found and connected to, use the Ping method:
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't Ping the database!", err)
	}
	fmt.Println("Pinged mongodb database successfully!")

	// To insert a document into a collection, first retrieve a Database and then Collection instance from the Client
	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Blog Service started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	fmt.Println("listenning on port 50051")

	// Better implementation with closing of server & listener upon "Ctrl+ c"

	// implementing ssl
	opts := []grpc.ServerOption{}
	certfile := "ssl/server.crt"
	keyfile := "ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certfile, keyfile)
	if sslErr != nil {
		log.Fatalf("failed loading certificates: %v", sslErr)
		return
	}

	opts = append(opts, grpc.Creds(creds))
	// s := grpc.NewServer(opts...)

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})
	// Register a reflection service on gRPC server
	//  (here, we are registering evans cli to interact with our gRPC service)
	reflection.Register(s)

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
	fmt.Println("Disconnecting the MongoDB connection")
	fmt.Println("End of Program")
}
