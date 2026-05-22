package main

import pb "hospital/generated/proto"

type server struct {
	pb.UnimplementedHospitalServiceServer
	store       *store
	consulStore *consultationStore
}
