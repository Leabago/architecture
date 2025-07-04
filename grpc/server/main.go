package main

import (
	"context"
	"flag"
	"fmt"
	"grpc/gen/go/hello"
	"grpc/gen/go/user"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	hello.UnimplementedGreetingsServer
}

type UserServer struct {
	user.UnimplementedUserServiceServer
}

func (s *server) SayHello(_ context.Context, in *hello.MessageRequest) (*hello.MessageResponse, error) {
	log.Printf("Received: %v", in.GetName())

	return &hello.MessageResponse{
		Message: "hello " + in.GetName(),
	}, nil
}

func (u *UserServer) UanryCall(_ context.Context, in *user.User) (*user.User, error) {
	log.Printf("UanryCall, Received: %v", in)
	in.Email = in.Email + "@mail.com"

	// var err error = fmt.Errorf("kek")

	return in, nil
}

func (u *UserServer) ServerStreammingCall(in *user.User, stream user.UserService_ServerStreammingCallServer) error {
	log.Printf("ServerStreammingCall: %s", in)

	in.Email += "@mail.com"

	for i := 0; i < 10; i++ {
		err := stream.Send(in)
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 2)
	}

	return nil
}

func (u *UserServer) ClientStreamingCall(stream user.UserService_ClientStreamingCallServer) error {

	lastUser := &user.User{}
	for {
		userResp, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				lastUser.Email += "@mail.com"
				return stream.SendAndClose(lastUser)
			}

			return err
		}

		fmt.Printf("ClientStreamingCall: %s \n", userResp)
		lastUser = userResp
	}
}

func (u *UserServer) BidirectionalStreaming(stream user.UserService_BidirectionalStreamingServer) error {
	lastUser := &user.User{Username: "hello"}
	fmt.Println(lastUser)

	for {
		userResp, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				break
			}
			fmt.Println("err 1: ", err)
			return err
		}

		userResp.Email += "@mail.com"
		if err := stream.Send(userResp); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.RegisterGreetingsServer(s, &server{})
	user.RegisterUserServiceServer(s, &UserServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
