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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName2 = "world"
	defaultAge2  = 10
)

var (
	addr2 = flag.String("addr", "localhost:50051", "the address to connect to")
	name2 = flag.String("name", defaultName2, "Name to greet")
	age2  = flag.Int64("age", defaultAge2, "Age to greet")
)

func main() {
	flag.Parse()
	// 建立一个与服务器的链接
	conn, err := grpc.Dial(*addr2, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 链接服务器并打印服务器的响应
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name2, Age: *age2})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s, age: %d", r.GetMessage(), r.GetAge())

	r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: *name2, Age: *age2})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s, Greetint Age: %d", r.GetMessage(), r.GetAge())
}
