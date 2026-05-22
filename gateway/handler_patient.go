package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	pb "hospital/generated/proto"
)

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
