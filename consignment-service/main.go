package main

import (
	"log"
	"net"

	pb "github.com/mygetzu/microservice-go-with-grpc-communication/consignment-service/proto/consignment"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

type service struct {
	repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		Created:     true,
		Consignment: consignment}, nil
}

func main() {
	repo := Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	} else {
		log.Printf("Listening to the port %v", port)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &service{&repo})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}
	log.Printf("gRPC service is running ..")
}
