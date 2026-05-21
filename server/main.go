package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"os"

	pb "hospital/generated/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedHospitalServiceServer
	store *store
}

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

func (s *server) ScheduleAppointment(_ context.Context, req *pb.ScheduleAppointmentRequest) (*pb.Appointment, error) {
	if req.ClinicId == 0 || req.DoctorId == 0 || req.PatientId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "clinic_id, doctor_id, and patient_id are required")
	}
	if req.ScheduledAt == "" {
		return nil, status.Errorf(codes.InvalidArgument, "scheduled_at is required")
	}

	dConflict, err := s.store.doctorConflict(req.DoctorId, req.ScheduledAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	if dConflict {
		return nil, status.Errorf(codes.AlreadyExists, "doctor already has an appointment at that time")
	}

	pConflict, err := s.store.patientConflict(req.PatientId, req.ScheduledAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	if pConflict {
		return nil, status.Errorf(codes.AlreadyExists, "patient already has an appointment at that time")
	}

	id, err := s.store.insertAppointment(req.ClinicId, req.DoctorId, req.PatientId, req.ScheduledAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	return &pb.Appointment{
		Id:          id,
		ClinicId:    req.ClinicId,
		DoctorId:    req.DoctorId,
		PatientId:   req.PatientId,
		ScheduledAt: req.ScheduledAt,
		Status:      "SCHEDULED",
	}, nil
}

func (s *server) ListAppointments(_ context.Context, req *pb.ListAppointmentsRequest) (*pb.ListAppointmentsResponse, error) {
	appts, err := s.store.listAppointments(req.ClinicId, req.DoctorId, req.PatientId, req.Date)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	resp := &pb.ListAppointmentsResponse{}
	for _, a := range appts {
		resp.Appointments = append(resp.Appointments, &pb.Appointment{
			Id:          a.id,
			ClinicId:    a.clinicID,
			DoctorId:    a.doctorID,
			PatientId:   a.patientID,
			ScheduledAt: a.scheduledAt,
			Status:      a.status,
		})
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

func (s *server) UpdateAppointmentStatus(_ context.Context, req *pb.UpdateAppointmentStatusRequest) (*pb.Appointment, error) {
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	switch req.Status {
	case "COMPLETED", "CANCELLED":
	default:
		return nil, status.Errorf(codes.InvalidArgument, "status must be COMPLETED or CANCELLED")
	}
	a, err := s.store.updateAppointmentStatus(req.Id, req.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "appointment not found")
		}
		return nil, status.Errorf(codes.Internal, "db error")
	}
	return &pb.Appointment{
		Id:          a.id,
		ClinicId:    a.clinicID,
		DoctorId:    a.doctorID,
		PatientId:   a.patientID,
		ScheduledAt: a.scheduledAt,
		Status:      a.status,
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
