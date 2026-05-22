package main

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type consultationDoc struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	AppointmentID int64         `bson:"appointment_id"`
	DoctorID      int64         `bson:"doctor_id"`
	PatientID     int64         `bson:"patient_id"`
	ClinicID      int64         `bson:"clinic_id"`
	Symptoms      string        `bson:"symptoms"`
	Prescription  string        `bson:"prescription"`
	Notes         string        `bson:"notes"`
	CreatedAt     time.Time     `bson:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at"`
}

type consultationStore struct {
	col *mongo.Collection
}

func newConsultationStore(uri string) (*consultationStore, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	col := client.Database("hospital_db").Collection("consultations")

	_, _ = col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "appointment_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	return &consultationStore{col: col}, nil
}

func (cs *consultationStore) upsert(ctx context.Context, doc consultationDoc) (*consultationDoc, error) {
	doc.UpdatedAt = time.Now()

	filter := bson.M{"appointment_id": doc.AppointmentID}
	update := bson.M{
		"$set": bson.M{
			"doctor_id":    doc.DoctorID,
			"patient_id":   doc.PatientID,
			"clinic_id":    doc.ClinicID,
			"symptoms":     doc.Symptoms,
			"prescription": doc.Prescription,
			"notes":        doc.Notes,
			"updated_at":   doc.UpdatedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result consultationDoc
	err := cs.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cs *consultationStore) findByAppointmentID(ctx context.Context, appointmentID int64) (*consultationDoc, error) {
	var doc consultationDoc
	err := cs.col.FindOne(ctx, bson.M{"appointment_id": appointmentID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &doc, nil
}
