to generate pb file
protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

1.to run server:-
go run greet_server/server.go

2.to run client:-
go run greet_client/client.go
  
#created a rpc for doing
1. Unary Streaming
2. Server Streaming
3. Client Streaming
4. Bidirectional Streaming
5. Unart Streaming with deadline
