package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"

	pb "soln/src/proto"

	"time"

	"google.golang.org/grpc"

	cmn "soln/src/common"
)

func main() {
	rand.Seed(time.Now().Unix())

	// dail server
	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	client := pb.NewMathClient(conn)
	stream, err := client.Max(context.Background())
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	var max uint32
	ctx := stream.Context()
	done := make(chan bool)

	// first goroutine sends random increasing numbers to stream
	// and closes int after 10 iterations
	go func() {
		for i := 1; i <= 10; i++ {
			// generate random nummber and send it to stream
			rnd := uint32(rand.Intn(i * 100))

			signer, err := cmn.LoadPrivateKey("private.pem")
			if err != nil {
				fmt.Errorf("signer is damaged: %v", err)
			}

			//toSign := "date: Thu, 05 Jan 2012 21:31:40 GMT"
			str := fmt.Sprint(rnd)
			fmt.Printf("String %v\n", str)

			signed, err := signer.Sign([]byte(str))
			if err != nil {
				fmt.Errorf("could not sign request: %v", err)
			}
			sig := base64.StdEncoding.EncodeToString(signed)

			//req := pb.Request{Num: rnd}
			req := pb.Request{Num: rnd, Msg: sig}

			if err := stream.Send(&req); err != nil {
				log.Fatalf("can not send %v", err)
			}
			log.Printf("%d sent", req.Msg)
			time.Sleep(time.Millisecond * 200)
		}
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}()

	// second goroutine receives data from stream
	// and saves result in max variable
	//
	// if stream is finished it closes done channel
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			max = resp.Result
			log.Printf("new max %d received", max)
		}
	}()

	// third goroutine closes done channel
	// if context is done
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	<-done
	log.Printf("finished with max=%d", max)
}
