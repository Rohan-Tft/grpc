package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/rohan/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("hello")

	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatal("failed to connect : ", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	fmt.Printf("created the client : %f", c)

	 doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doBidirectionalStreaming(c)

	doUnaryWithDeadline(c, 5*time.Second) //should complete
	doUnaryWithDeadline(c, 1*time.Second) //should timeout

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting server streaming ")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rohan",
			LastName:  "Yadav",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling grpc greetmany times %v", err)
	}
	for {
		mes, err := resStream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from greetmany times is %v", mes.GetResult())

	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting Do Unary  RPC ")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rohan",
			LastName:  "Yadav",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling greet rpc: %v", err)
	}
	log.Printf("Response from greet is :  %v", res)

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting client streaming  RPC ")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "rohan",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "jayant",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "shyam",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "mohit",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling longgreet %v", err)
	}
	//iterating over our slice and sending the each message individually
	for _, req := range requests {
		fmt.Printf("sending req %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	fmt.Println()

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recienving response %v", err)
	}
	log.Printf("longgreet response %v\n", res)

}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Bi directional streaming client RPC")

	//4 steps
	//1--create a stream by invoking the client
	//2--we send bunch of messages
	//3--we recienve bunch of messages
	//block until everything is done

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("Error while creating stream %v", err)
	}

	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "rohan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "jayant",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "shyam",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "mohit",
			},
		},
	}
	//creating channel for blocking untill everthing is done

	waitc := make(chan struct{})

	// Send bunch of messages
	go func() {
		for _, req := range requests {

			fmt.Printf("Sending message request %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()

	}()

	go func() {
		//Recieving messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while recieving stream%v", err)
				break
			}
			fmt.Println("recieved message : ", res.GetResult())

		}
		close(waitc)
	}()
	//blocking
	<-waitc

}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, seconds time.Duration) {
	fmt.Println("starting Do UnaryWithDeadline  RPC ")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rohan",
			LastName:  "Yadav",
		},
	}
	ctx, cancle := context.WithTimeout(context.Background(), seconds)
	defer cancle()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		StatusErr, ok := status.FromError(err)

		if ok {
			if StatusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout ! deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error %v", StatusErr)
			}
		} else {

			log.Fatalf("error while calling greetWithDeadline rpc: %v", err)
		}
		return
	}
	log.Printf("Response from greetWithDeadline is :  %v", res)

}
