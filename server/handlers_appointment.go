package main

import (
	"context"
	"database/sql"
	"errors"

	pb "hospital/generated/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

	pending, err := s.store.doctorHasPendingAppointment(req.DoctorId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "db error")
	}
	if pending {
		return nil, status.Errorf(codes.FailedPrecondition, "doctor has a pending appointment — mark it completed first")
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
