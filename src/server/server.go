package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"

	pb "soln/src/proto"

	"google.golang.org/grpc"

	cmn "soln/src/common"
)

type server struct{}

func (s server) Max(srv pb.Math_MaxServer) error {

	log.Println("start new server")
	var max uint32
	ctx := srv.Context()

	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		parser, perr := cmn.LoadPublicKey("public.pem")
		if perr != nil {
			fmt.Errorf("could not sign request: %v", err)
		}
		signed_msg := []byte(req.Msg)

		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, req.Num)

		err = parser.Unsign(bs, signed_msg)
		if err != nil {
			fmt.Errorf("could not sign request: %v", err)
		}

		// continue if number reveived from stream
		// less than max
		cNum := uint32(req.Num)

		if cNum <= max {
			continue
		}

		// update max and send it to stream
		max = cNum
		resp := pb.Response{Result: max}
		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Printf("send new max=%d", max)
	}
}

func main() {
	// create listiner
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	s := grpc.NewServer()
	pb.RegisterMathServer(s, server{})

	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
