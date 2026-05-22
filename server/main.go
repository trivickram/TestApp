package main

import (
	"log"
	"net"
	"os"

	pb "hospital/generated/proto"

	"google.golang.org/grpc"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/hospital?parseTime=true"
	}

	st, err := newStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	cs, err := newConsultationStore(mongoURI)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	log.Printf("mongodb connected at %s", mongoURI)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHospitalServiceServer(grpcServer, &server{store: st, consulStore: cs})

	log.Println("grpc server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
