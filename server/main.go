package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "hospital/generated/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedHospitalServiceServer
	store *store
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

	if err := s.store.insertPatient(patient.Id, patient.Name, patient.Age); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save patient")
	}

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

	exists, err := s.store.patientExists(req.GetPatientId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "patient not found")
	}

	appointment := &pb.Appointment{
		Id:          uuid.NewString(),
		PatientId:   req.GetPatientId(),
		Doctor:      req.GetDoctor(),
		ScheduledAt: req.GetScheduledAt(),
		Status:      "SCHEDULED",
	}

	if err := s.store.insertAppointment(appointment.Id, appointment.PatientId, appointment.Doctor, appointment.ScheduledAt, appointment.Status); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save appointment")
	}

	return appointment, nil
}

func (s *server) GetAppointmentStatus(_ context.Context, req *pb.GetAppointmentStatusRequest) (*pb.AppointmentStatus, error) {
	if req.GetAppointmentId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "appointment_id is required")
	}

	st, found, err := s.store.getAppointmentStatus(req.GetAppointmentId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "appointment not found")
	}

	return &pb.AppointmentStatus{
		AppointmentId: req.GetAppointmentId(),
		Status:        st,
	}, nil
}

func (s *server) UpdateAppointmentStatus(_ context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.AppointmentStatus, error) {
	if req.GetAppointmentId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "appointment_id is required")
	}
	if req.GetStatus() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "status is required")
	}

	updated, err := s.store.updateAppointmentStatus(req.GetAppointmentId(), req.GetStatus())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	if !updated {
		return nil, status.Errorf(codes.NotFound, "appointment not found")
	}

	return &pb.AppointmentStatus{
		AppointmentId: req.GetAppointmentId(),
		Status:        req.GetStatus(),
	}, nil
}

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/hospital?parseTime=true"
	}

	st, err := newStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHospitalServiceServer(grpcServer, &server{store: st})

	log.Println("grpc server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
