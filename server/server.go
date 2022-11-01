package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"grpc_test1/model"
	pb "grpc_test1/pb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

type Server struct {
    pb.UnimplementedServiceServerServer
	// savedReply []*pb.GetNameReply // read-only after initialized
}

func (s *Server) Getname(ctx context.Context, in *pb.GetIdRequest) (*pb.GetNameReply, error) {
	tempid, err := strconv.Atoi(in.GetDbId())
	fmt.Println(int(tempid))
	name, err := model.QueryNameWithID(int(tempid))
	if err != nil {
		fmt.Printf("Cant find this id, return default")
		name = "default"
	}

	log.Printf("Received: %v", in.GetDbId() )
	return &pb.GetNameReply{Name: name}, nil
}

func (s *Server) OGetMname( in *pb.GetIdRequest, stream pb.ServiceServer_OGetMnameServer ) error {
	var temp pb.GetNameReply
	if in.DbId == "1" {
		temp.Name = "Jerry88819"
	} else {
		temp.Name = "qqqqqqqqqq"
	}
	
	for i := 0; i < 10; i++ {
		if err := stream.Send(&temp); err != nil {
			return err
		}
	}

	return nil
}

func( s *Server ) MGetOname( stream pb.ServiceServer_MGetOnameServer ) error {
	var temp string
	for {
		
		request, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("right")
			return stream.SendAndClose(&pb.GetNameReply{Name: temp})
		}

		temp = temp + request.DbId
		fmt.Println(temp)

		if err != nil {
			fmt.Println("qweqweq")
			return err
			// return stream.SendAndClose(&pb.GetNameReply{Name: temp})
		}
		// temp = temp + request.DbId
		// fmt.Println(temp)
		
	}
} // MGetOname()

func ( s *Server ) MGetMname( stream pb.ServiceServer_MGetMnameServer ) error {
	for {
		receive, err := stream.Recv()
		if err == io.EOF { // 訊息傳遞完畢 就結束了
			return nil
		}
		if err != nil {
			log.Fatal("Receive message error!")
			return err
		}

		var temp string
		if receive.DbId == "1" {
			temp = "one"
		} else if receive.DbId == "2" {
			temp = "two"
		} else {
			temp = "wow"
		}

		
		if err := stream.Send(&pb.GetNameReply{Name:temp}); err != nil {
			return err
		}
	}
} // MGetMname()

func main() {

    db, err := sql.Open("sqlite3", "./foo.db")
    if err != nil {
		fmt.Printf("database error =>%v\n", err)
		return
	}
    fmt.Printf("database init success\n")

    model.Init(db)
    defer db.Close()


    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 3000))
    if err != nil {
		log.Fatalf("failed to listed: %v", err)
	}
    checkErr(err)

    s := grpc.NewServer()
	pb.RegisterServiceServerServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

    return
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}