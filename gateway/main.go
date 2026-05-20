package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "hospital/generated/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type gateway struct {
	client pb.HospitalServiceClient
}

type errBody struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func grpcError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, errBody{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusBadRequest, errBody{Error: st.Message()})
}

func (g *gateway) createPatient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errBody{Error: "method not allowed"})
		return
	}

	var req pb.CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errBody{Error: "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := g.client.CreatePatient(ctx, &req)
	if err != nil {
		grpcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (g *gateway) scheduleAppointment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errBody{Error: "method not allowed"})
		return
	}

	var req pb.ScheduleAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errBody{Error: "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := g.client.ScheduleAppointment(ctx, &req)
	if err != nil {
		grpcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (g *gateway) getAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, errBody{Error: "method not allowed"})
		return
	}

	appointmentID := r.URL.Query().Get("appointment_id")
	if appointmentID == "" {
		writeJSON(w, http.StatusBadRequest, errBody{Error: "appointment_id is required"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := g.client.GetAppointmentStatus(ctx, &pb.GetAppointmentStatusRequest{AppointmentId: appointmentID})
	if err != nil {
		grpcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (g *gateway) updateAppointmentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errBody{Error: "method not allowed"})
		return
	}

	var req pb.UpdateAppointmentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errBody{Error: "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := g.client.UpdateAppointmentStatus(ctx, &req)
	if err != nil {
		grpcError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	g := &gateway{client: pb.NewHospitalServiceClient(conn)}

	mux := http.NewServeMux()
	mux.HandleFunc("/patients", g.createPatient)
	mux.HandleFunc("/appointments", g.scheduleAppointment)
	mux.HandleFunc("/appointments/status", g.getAppointmentStatus)
	mux.HandleFunc("/appointments/status/update", g.updateAppointmentStatus)

	log.Println("gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}
