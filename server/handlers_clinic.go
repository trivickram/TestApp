package main

import (
	"context"

	pb "hospital/generated/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) ListClinics(_ context.Context, _ *pb.Empty) (*pb.ListClinicsResponse, error) {
	clinics, err := s.store.listClinics()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	resp := &pb.ListClinicsResponse{}
	for _, c := range clinics {
		resp.Clinics = append(resp.Clinics, &pb.Clinic{Id: c.id, Name: c.name})
	}
	return resp, nil
}
