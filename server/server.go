package main

import (
	pb "counter/counter"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := 8070
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	countServer := &CountServer{}
	pb.RegisterCounterServer(grpcServer, countServer)

	log.Printf("starting server at 0.0.0.0:%d", port)
	grpcServer.Serve(lis)
}

type CountServer struct{}

func (cs *CountServer) CountValue(stream pb.Counter_CountValueServer) error {

	msgCount := 1
	for {
		count, err := stream.Recv()

		if err != nil {
			return err
		}
		log.Println("received from client:", count.Value)
		if msgCount == 3 {
			return stream.SendAndClose(&pb.Error{Msg: "closing stream"})
		}
		msgCount++

	}

}
