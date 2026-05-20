package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "hospital/generated/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedHospitalServiceServer
	mu           sync.RWMutex
	patients     map[string]*pb.Patient
	appointments map[string]*pb.Appointment
}

func newServer() *server {
	return &server{
		patients:     map[string]*pb.Patient{},
		appointments: map[string]*pb.Appointment{},
	}
}

func (s *server) CreatePatient(_ context.Context, req *pb.CreatePatientRequest) (*pb.Patient, error) {
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.GetAge() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "age must be greater than 0")
	}

	patient := &pb.Patient{
		Id:   uuid.NewString(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	s.mu.Lock()
	s.patients[patient.GetId()] = patient
	s.mu.Unlock()

	return patient, nil
}

func (s *server) ScheduleAppointment(_ context.Context, req *pb.ScheduleAppointmentRequest) (*pb.Appointment, error) {
	if req.GetPatientId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "patient_id is required")
	}
	if req.GetDoctor() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "doctor is required")
	}
	if req.GetScheduledAt() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "scheduled_at is required")
	}

	s.mu.RLock()
	_, ok := s.patients[req.GetPatientId()]
	s.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "patient not found")
	}

	appointment := &pb.Appointment{
		Id:          uuid.NewString(),
		PatientId:   req.GetPatientId(),
		Doctor:      req.GetDoctor(),
		ScheduledAt: req.GetScheduledAt(),
		Status:      "SCHEDULED",
	}

	s.mu.Lock()
	s.appointments[appointment.GetId()] = appointment
	s.mu.Unlock()

	return appointment, nil
}

func (s *server) GetAppointmentStatus(_ context.Context, req *pb.GetAppointmentStatusRequest) (*pb.AppointmentStatus, error) {
	if req.GetAppointmentId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "appointment_id is required")
	}

	s.mu.RLock()
	appointment, ok := s.appointments[req.GetAppointmentId()]
	s.mu.RUnlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "appointment not found")
	}

	return &pb.AppointmentStatus{
		AppointmentId: appointment.GetId(),
		Status:        appointment.GetStatus(),
	}, nil
}

func (s *server) UpdateAppointmentStatus(_ context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.AppointmentStatus, error) {
	if req.GetAppointmentId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "appointment_id is required")
	}
	if req.GetStatus() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "status is required")
	}

	s.mu.Lock()
	appointment, ok := s.appointments[req.GetAppointmentId()]
	if !ok {
		s.mu.Unlock()
		return nil, status.Errorf(codes.NotFound, "appointment not found")
	}
	appointment.Status = req.GetStatus()
	s.mu.Unlock()

	return &pb.AppointmentStatus{
		AppointmentId: appointment.GetId(),
		Status:        appointment.GetStatus(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHospitalServiceServer(grpcServer, newServer())

	log.Println("grpc server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
