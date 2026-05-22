package main

import (
	"context"

	pb "hospital/generated/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) CreateDoctor(_ context.Context, req *pb.CreateDoctorRequest) (*pb.Doctor, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	id, err := s.store.insertDoctor(req.Name, req.Specialization)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	return &pb.Doctor{Id: id, Name: req.Name, Specialization: req.Specialization}, nil
}

func (s *server) LinkDoctorToClinic(_ context.Context, req *pb.LinkDoctorRequest) (*pb.Empty, error) {
	if req.ClinicId == 0 || req.DoctorId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "clinic_id and doctor_id are required")
	}
	if err := s.store.linkDoctor(req.ClinicId, req.DoctorId); err != nil {
		if err == errAlreadyLinked {
			return nil, status.Errorf(codes.AlreadyExists, "doctor already linked to this clinic")
		}
		return nil, status.Errorf(codes.Internal, "db error")
	}
	return &pb.Empty{}, nil
}

func (s *server) ListClinicDoctors(_ context.Context, req *pb.ListClinicDoctorsRequest) (*pb.ListDoctorsResponse, error) {
	doctors, err := s.store.listClinicDoctors(req.ClinicId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	resp := &pb.ListDoctorsResponse{}
	for _, d := range doctors {
		resp.Doctors = append(resp.Doctors, &pb.Doctor{Id: d.id, Name: d.name, Specialization: d.specialization})
	}
	return resp, nil
}

func (s *server) SearchDoctors(_ context.Context, req *pb.SearchDoctorsRequest) (*pb.ListDoctorsResponse, error) {
	doctors, err := s.store.searchDoctors(req.Query, req.ClinicId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	resp := &pb.ListDoctorsResponse{}
	for _, d := range doctors {
		resp.Doctors = append(resp.Doctors, &pb.Doctor{Id: d.id, Name: d.name, Specialization: d.specialization})
	}
	return resp, nil
}
