package main

import (
	"context"
	"flag"
	"fmt"
	"grpc/gen/go/hello"
	"grpc/gen/go/user"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func SayHello(c hello.GreetingsClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := (c).SayHello(ctx, &hello.MessageRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

func UanryCall(us user.UserServiceClient) {

	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	user := &user.User{
		Username: "Vasya",
		Email:    "Vasya",
	}

	userResp, err := (us).UanryCall(cxt, user)
	if err != nil {
		log.Fatalf("could not sendUser: %v", err)
	}

	log.Printf("UanryCall:  %s", userResp)
}

func ServerStreammingCall(us user.UserServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	user := &user.User{
		Username: "Vasya",
		Email:    "Vasya",
	}

	stream, err := (us).ServerStreammingCall(ctx, user)
	if err != nil {
		log.Fatalf("could not sendUser: %v", err)
	}

	for {
		userResp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("client.ServerStreammingCall failed: %v", err)
		}

		fmt.Printf("ServerStreammingCall: %s \n", userResp)
	}

	time.Sleep(time.Second * 5)

	return nil
}

func ClientStreamingCall(us user.UserServiceClient) error {
	cxt, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := (us).ClientStreamingCall(cxt)
	if err != nil {
		log.Fatalf("could create stream: %v", err)
	}

	user := &user.User{
		Username: "Vasya",
		Email:    "Vasya",
	}

	for i := 0; i < 10; i++ {
		err := stream.Send(user)
		if err != nil {
			fmt.Println("err: ", err)
			return err
		}
	}

	userResp, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	fmt.Printf("ClientStreamingCall: %s \n", userResp)
	return nil
}

func BidirectionalStreaming(us user.UserServiceClient) error {

	waitc := make(chan struct{})

	stream, err := (us).BidirectionalStreaming(context.Background())
	if err != nil {
		fmt.Println("err: ", err)
		return err
	}

	u := &user.User{
		Username: "Vasya",
		Email:    "Vasya",
	}

	go func() {
		for {
			userRes := &user.User{}
			userRes, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					close(waitc)
					break
				}
				fmt.Println("err: ", err)
				log.Fatal(err)
			}
			fmt.Println("userRes: ", userRes)
		}
	}()

	// go func() {
	for i := 0; i < 10; i++ {
		fmt.Println("Send")
		err := stream.Send(u)
		if err != nil {
			fmt.Println("err: ", err)
			log.Fatal(err)
		}
	}

	stream.CloseSend()
	<-waitc
	return nil
}

func main() {

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Greetings
	client := hello.NewGreetingsClient(conn)

	// Contact the server and print out its response.
	SayHello(client)

	// User
	userClient := user.NewUserServiceClient(conn)
	UanryCall(userClient)
	ServerStreammingCall(userClient)
	ClientStreamingCall(userClient)
	BidirectionalStreaming(userClient)
}
