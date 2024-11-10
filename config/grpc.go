package config

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RpcDial(targetPort string) (*grpc.ClientConn, error) {
	grpcConn, err := grpc.Dial(
		"localhost:"+targetPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*64),
			grpc.MaxCallSendMsgSize(1024*1024*64),
		),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	log.Printf("[Running-Success] gRPC clients connection")
	return grpcConn, nil
}
