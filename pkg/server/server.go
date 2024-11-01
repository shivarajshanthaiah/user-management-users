package server

import (
	"fmt"
	"log"
	"net"

	"github.com/shivaraj-shanthaiah/user-management/pkg/handler"
	pb "github.com/shivaraj-shanthaiah/user-management/pkg/proto"
	"google.golang.org/grpc"
)

func NewGrpcUserServer(port string, handler *handler.UserHandler) error {
	log.Println("Starting gRPC server")

	addr := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("error creating listener on %v", addr)
		return err
	}

	grpc := grpc.NewServer()
	pb.RegisterUserServiceServer(grpc, handler)

	log.Printf("listening on gRPC server at: %v", port)
	err = grpc.Serve(lis)
	if err != nil {
		log.Printf("Error serving gRPC server at %v: %s", addr, err.Error())
		return err
	}
	return nil
}
