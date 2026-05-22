package main

import (
	"context"
	"errors"
	"time"

	pb "hospital/generated/proto"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) SaveConsultation(_ context.Context, req *pb.SaveConsultationRequest) (*pb.Consultation, error) {
	if req.AppointmentId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "appointment_id is required")
	}
	if req.Symptoms == "" {
		return nil, status.Errorf(codes.InvalidArgument, "symptoms are required")
	}
	if req.Prescription == "" {
		return nil, status.Errorf(codes.InvalidArgument, "prescription is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := s.consulStore.upsert(ctx, consultationDoc{
		AppointmentID: req.AppointmentId,
		DoctorID:      req.DoctorId,
		PatientID:     req.PatientId,
		ClinicID:      req.ClinicId,
		Symptoms:      req.Symptoms,
		Prescription:  req.Prescription,
		Notes:         req.Notes,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save consultation")
	}
	return docToProto(doc), nil
}

func (s *server) GetConsultation(_ context.Context, req *pb.GetConsultationRequest) (*pb.Consultation, error) {
	if req.AppointmentId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "appointment_id is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := s.consulStore.findByAppointmentID(ctx, req.AppointmentId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "no consultation for this appointment")
		}
		return nil, status.Errorf(codes.Internal, "db error")
	}
	return docToProto(doc), nil
}

func docToProto(doc *consultationDoc) *pb.Consultation {
	return &pb.Consultation{
		Id:            doc.ID.Hex(),
		AppointmentId: doc.AppointmentID,
		DoctorId:      doc.DoctorID,
		PatientId:     doc.PatientID,
		ClinicId:      doc.ClinicID,
		Symptoms:      doc.Symptoms,
		Prescription:  doc.Prescription,
		Notes:         doc.Notes,
		CreatedAt:     doc.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     doc.UpdatedAt.Format(time.RFC3339),
	}
}
