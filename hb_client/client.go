package main

import (
	"context"
	"fmt"
	heartbeat_pb "bmutziu.me/hb_proto"
	"io"
	"log"
	"math/rand"
	"sync"

	"google.golang.org/grpc"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generateHeartBeat() int32 {
	bpm := rand.Intn(100)
	return int32(bpm)
}

var wg sync.WaitGroup

func NormalAbnormalHeartBeat(c heartbeat_pb.HeartBeatServiceClient) {
	stream, err := c.NormalAbnormalHeartBeat(context.Background())
	handleError(err)
	for t := 0; t < 4; t++ {
		newBpm := generateHeartBeat()

		newNormalAbnormalHeartBeatRequest := &heartbeat_pb.NormalAbnormalHeartBeatRequest{
			Heartbeat: &heartbeat_pb.NormalAbnormalHeartBeat{
				Bpm: newBpm,
			},
		}
		stream.Send(newNormalAbnormalHeartBeatRequest)
		fmt.Printf("Sent %v\n", newNormalAbnormalHeartBeatRequest)
	}
	stream.CloseSend()

	wg.Add(1)
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				wg.Done()
				break
			}
			handleError(err)
			fmt.Printf("Received %v\n", msg)
		}
	}()

	wg.Wait()
}

func HeartBeatHistory(c heartbeat_pb.HeartBeatServiceClient) {
	heartBeatHistoryRequest := &heartbeat_pb.HeartBeatHistoryRequest{
		Username: "bmutziulhb",
	}
	res_stream, _ := c.HeartBeatHistory(context.Background(), heartBeatHistoryRequest)

	for {
		msg, err := res_stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println(msg)
	}

}

func LiveHeartBeat(c heartbeat_pb.HeartBeatServiceClient) {
	stream, err := c.LiveHeartBeat(context.Background())
	handleError(err)
	var username = "bmutziulhb"

	for t := 0; t < 8; t++ {
		newBpm := generateHeartBeat()
		newLiveHeartBeatRequest := &heartbeat_pb.LiveHeartBeatRequest{
			Heartbeat: &heartbeat_pb.HeartBeat{
				Bpm:      newBpm,
				Username: username,
			},
		}

		fmt.Println("Request Sent: ", newLiveHeartBeatRequest)
		stream.Send(newLiveHeartBeatRequest)
	}

	resp, err := stream.CloseAndRecv()
	handleError(err)
	fmt.Println(resp)
}

func HeartBeat(c heartbeat_pb.HeartBeatServiceClient) {
	heartbeatRequest := heartbeat_pb.HeartBeatRequest{
		Heartbeat: &heartbeat_pb.HeartBeat{
			Bpm:      73,
			Username: "bmutziu",
		},
	}

	res, err := c.UserHeartBeat(context.Background(), &heartbeatRequest)
	handleError(err)

	log.Printf("Response from server: %v", res)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	handleError(err)

	fmt.Println("Client Running")
	defer conn.Close()

	c := heartbeat_pb.NewHeartBeatServiceClient(conn)
	// HeartBeat(c)
	// LiveHeartBeat(c)
	// HeartBeatHistory(c)
	NormalAbnormalHeartBeat(c)
}
