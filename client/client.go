package main

import (
	"context"
	pb "counter/counter"
	"errors"
	"google.golang.org/grpc"
	"time"

	"log"
)

const ServerCloseErr = "server close"

func main() {

	serverAddr := "0.0.0.0:8070"

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewCounterClient(conn)

	c := &Client{count: 0}

	for {
		err = c.startCounter(client)
		if err.Error() != ServerCloseErr {
			break
		}
		log.Println("reconnecting")
	}

	log.Println("End of Client")
}

type Client struct {
	count int32
}

func (c *Client) startCounter(client pb.CounterClient) error {

	ticker := time.NewTicker(5 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.CountValue(ctx)

	if err != nil {
		log.Fatalf("%v", err)
	}

	errChan := make(chan error, 1)
	done := make(chan bool)
	go func() {
		msg := &pb.Error{}
		rErr := stream.RecvMsg(msg)
		log.Println("server sent:", msg.Msg)
		//c.count--
		if rErr == nil {
			rErr = errors.New("server close")
		}
		errChan <- rErr
		done <- true

		log.Println("ending receive")
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println("sending count", c.count)
				count := &pb.Count{Value: c.count}
				err := stream.Send(count)
				c.count++
				if err != nil {
					log.Printf("%v", err)
				}
			case <-done:
				return
			}
		}
	}()

	err = <-errChan
	ticker.Stop()

	return err
}
