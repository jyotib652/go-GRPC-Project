# Old
# protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

# New
# protoc --go-grpc_out=. path/to/file.proto

# protoc-gen-go-grpc is a plugin for the Google protocol buffer compiler to generate Go code
# Now "protoc-gen-go-grpc" package provides the code to generate .go files from .proto files


# that mean I'll need both protoc-gen-go and protoc-gen-go-grpc to generate my gRPC service definitions,
# and that protoc-gen-go is only deprecating the support for gRPC plugin and not messages?

# For greet
protoc --go-grpc_out=. greet/greetpb/greet.proto
protoc --go_out=plugins=. /greet/greetpb/greet.proto # support for grpc plugins gone but the code is supported so it becomes,
protoc --go_out=. greet/greetpb/greet.proto

# For calculator
protoc --go-grpc_out=. calculator/calculatorpb/calculator.proto
protoc --go_out=. calculator/calculatorpb/calculator.proto

# For blog
protoc --go-grpc_out=. blog/blogpb/blog.proto
protoc --go_out=. blog/blogpb/blog.proto

# protoc -I./protos        \
#   --go_out=./server      \
#   --go-grpc_out=./server \
#   --go_out=./client      \
#   --go-grpc_out=./client \
#   protos/*.proto 
# Reference: https://stackoverflow.com/questions/71777702/service-compiling-successfully-but-message-structs-not-generating-grpc-go