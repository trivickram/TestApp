package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	pb "hospital/generated/proto"
)

func (g *gateway) saveConsultation(w http.ResponseWriter, r *http.Request) {
	var req pb.SaveConsultationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.SaveConsultation(ctx, &req)
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (g *gateway) getConsultation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid appointment id"})
		return
	}
	ctx, cancel := rctx()
	defer cancel()
	resp, err := g.client.GetConsultation(ctx, &pb.GetConsultationRequest{AppointmentId: id})
	if err != nil {
		grpcErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
