package main

import (
	"context"
	"fmt"
	gtrace "github.com/moxiaomomo/grpc-jaeger"
	"google.golang.org/grpc"
	pb "grpc2/proto"
	"net"
	"os"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
func main() {
	var servOpts []grpc.ServerOption
	tracer, _, err := gtrace.NewJaegerTracer("testSrv", "192.168.129.130:6831")
	if err != nil {
		fmt.Printf("new tracer err: %+v\n", err)
		os.Exit(-1)
	}
	if tracer != nil {
		servOpts = append(servOpts, gtrace.ServerOption(tracer))
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer(servOpts...)
	pb.RegisterGreeterServer(s, &server{})
	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
}
