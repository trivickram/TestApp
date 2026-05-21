package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "hospital/generated/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type gateway struct {
	client pb.HospitalServiceClient
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func grpcErr(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	code := http.StatusInternalServerError
	switch st.Code() {
	case codes.InvalidArgument:
		code = http.StatusBadRequest
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.AlreadyExists:
		code = http.StatusConflict
	case codes.FailedPrecondition:
		code = http.StatusUnprocessableEntity
	}
	writeJSON(w, code, map[string]string{"error": st.Message()})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func rctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (g *gateway) listClinics(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.ListClinics(ctx, &pb.Empty{})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp.Clinics)
}

func (g *gateway) createPatient(w http.ResponseWriter, r *http.Request) {
	var req pb.CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.CreatePatient(ctx, &req)
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (g *gateway) createDoctor(w http.ResponseWriter, r *http.Request) {
	var req pb.CreateDoctorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.CreateDoctor(ctx, &req)
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (g *gateway) linkDoctor(w http.ResponseWriter, r *http.Request) {
	clinicID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid clinic id"})
		return
	}
	var body struct {
		DoctorId int64 `json:"doctor_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	_, err = g.client.LinkDoctorToClinic(ctx, &pb.LinkDoctorRequest{ClinicId: clinicID, DoctorId: body.DoctorId})
	if err != nil {
		grpcErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (g *gateway) listClinicDoctors(w http.ResponseWriter, r *http.Request) {
	clinicID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid clinic id"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.ListClinicDoctors(ctx, &pb.ListClinicDoctorsRequest{ClinicId: clinicID})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp.Doctors)
}

func (g *gateway) listClinicPatients(w http.ResponseWriter, r *http.Request) {
	clinicID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid clinic id"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.ListClinicPatients(ctx, &pb.ListClinicPatientsRequest{ClinicId: clinicID})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp.Patients)
}

func (g *gateway) scheduleAppointment(w http.ResponseWriter, r *http.Request) {
	var req pb.ScheduleAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.ScheduleAppointment(ctx, &req)
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (g *gateway) listAppointments(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	toInt := func(s string) int64 {
		v, _ := strconv.ParseInt(s, 10, 64)
		return v
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.ListAppointments(ctx, &pb.ListAppointmentsRequest{
		ClinicId:  toInt(q.Get("clinic_id")),
		DoctorId:  toInt(q.Get("doctor_id")),
		PatientId: toInt(q.Get("patient_id")),
		Date:      q.Get("date"),
	})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp.Appointments)
}

func (g *gateway) searchDoctors(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	clinicID, _ := strconv.ParseInt(r.URL.Query().Get("clinic_id"), 10, 64)
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.SearchDoctors(ctx, &pb.SearchDoctorsRequest{Query: q, ClinicId: clinicID})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp.Doctors)
}

func (g *gateway) searchPatients(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.SearchPatients(ctx, &pb.SearchPatientsRequest{Query: q})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp.Patients)
}

func (g *gateway) updateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid appointment id"})
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.UpdateAppointmentStatus(ctx, &pb.UpdateAppointmentStatusRequest{Id: id, Status: body.Status})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
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

	log.Println("gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}
