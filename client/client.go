package main

import (
	"context"
	"fmt"
	"io"

	// "io"
	"log"
	"strconv"
	"time"

	pb "grpc_test1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewServiceServerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// r, err := c.Getname(ctx, &pb.GetIdRequest{DbId: "2"})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }

	// log.Printf("Greeting: %s", r.GetName())

	// test2, err := c.OGetMname(ctx, &pb.GetIdRequest{DbId: "1"})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }

	// for {
	// 	feature, err := test2.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("client.ListFeatures failed: %v", err)
	// 	}
	// 	log.Printf("Feature: name: %s", feature.GetName())
	// }

	test3, err := c.MGetOname(ctx)
	if err != nil {
		log.Fatalf("client.M to O failed: %v", err)
	}

	var a []*pb.GetIdRequest
	for i := 0; i < 5 ; i++ {
		temp_string := strconv.Itoa(i)
		ttt := &pb.GetIdRequest{DbId: temp_string}
		a = append(a, ttt)
	}

	for _, i := range(a) {
		if err := test3.Send(i); err != nil {
			log.Fatalf("client.RecordRoute: stream.Send(%v) failed: %v", i, err)
		}
	}

	fmt.Println("asdasdasd")

	reply, err  := test3.CloseAndRecv()

	if err != nil {
		fmt.Println("hii")
		log.Fatalf("client.RecordRoute failed: %v", err)
	}

	log.Println("Get Name reply:", reply.Name)




	test4, err := c.MGetMname(ctx) 
	waitc := make(chan struct{})
	if err != nil {
		log.Fatalf("client.M to M failed: %v", err)
	}

	go func() {
		for {
			in, err := test4.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
			}

			fmt.Println(in.Name)
			
		}
	}()
	for _, note := range a {
		if err := test4.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", note, err)
		}
	}

	test4.CloseSend()
	<-waitc

}