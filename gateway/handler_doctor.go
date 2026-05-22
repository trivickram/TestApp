package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	pb "hospital/generated/proto"
)

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
