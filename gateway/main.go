package main

import (
	"log"
	"net/http"

	pb "hospital/generated/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type gateway struct {
	client pb.HospitalServiceClient
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	g := &gateway{client: pb.NewHospitalServiceClient(conn)}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /clinics", g.listClinics)
	mux.HandleFunc("POST /patients", g.createPatient)
	mux.HandleFunc("POST /doctors", g.createDoctor)
	mux.HandleFunc("POST /clinics/{id}/doctors", g.linkDoctor)
	mux.HandleFunc("GET /clinics/{id}/doctors", g.listClinicDoctors)
	mux.HandleFunc("GET /clinics/{id}/patients", g.listClinicPatients)
	mux.HandleFunc("POST /appointments", g.scheduleAppointment)
	mux.HandleFunc("GET /appointments", g.listAppointments)
	mux.HandleFunc("GET /doctors", g.searchDoctors)
	mux.HandleFunc("GET /patients", g.searchPatients)
	mux.HandleFunc("PATCH /appointments/{id}/status", g.updateAppointmentStatus)
	mux.HandleFunc("POST /consultations", g.saveConsultation)
	mux.HandleFunc("GET /appointments/{id}/consultation", g.getConsultation)

	log.Println("gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}
