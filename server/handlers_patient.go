package main

import (
	"context"

	pb "hospital/generated/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreatePatient(_ context.Context, req *pb.CreatePatientRequest) (*pb.Patient, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.Age <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "age must be positive")
	}
	id, err := s.store.insertPatient(req.Name, req.Age)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	return &pb.Patient{Id: id, Name: req.Name, Age: req.Age}, nil
}

func (s *server) ListClinicPatients(_ context.Context, _ *pb.ListClinicPatientsRequest) (*pb.ListPatientsResponse, error) {
	patients, err := s.store.listPatients()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	resp := &pb.ListPatientsResponse{}
	for _, p := range patients {
		resp.Patients = append(resp.Patients, &pb.Patient{Id: p.id, Name: p.name, Age: p.age})
	}
	return resp, nil
}

func (s *server) SearchPatients(_ context.Context, req *pb.SearchPatientsRequest) (*pb.ListPatientsResponse, error) {
	if req.Query == "" {
		return &pb.ListPatientsResponse{}, nil
	}
	patients, err := s.store.searchPatients(req.Query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	resp := &pb.ListPatientsResponse{}
	for _, p := range patients {
		resp.Patients = append(resp.Patients, &pb.Patient{Id: p.id, Name: p.name, Age: p.age})
	}
	return resp, nil
}
