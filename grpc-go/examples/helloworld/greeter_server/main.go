/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v, Age: %v", in.GetName(), in.GetAge())
	return &pb.HelloReply{Message: "Hello " + in.GetName(), Age: in.GetAge()}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("ReceivedAgent: %v, Age: %v", in.GetName(), in.GetAge())
	return &pb.HelloReply{Message: "Hello again " + in.GetName(), Age: in.GetAge()}, nil
}

func (s *server) PayIn(ctx context.Context, in *pb.PayInRequest) (*pb.PayInResponse, error) {
	log.Printf("Received PayIn provider : %v, amount: %v, phone: %v, appid: %v", in.GetProvider(), in.GetAmount(), in.GetPhone(), in.GetAppid())
	return &pb.PayInResponse{Sn: "Hello again " + in.GetProvider(), Appid: in.GetAppid(), Provider: in.GetProvider(), CreateTime: in.Amount}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
