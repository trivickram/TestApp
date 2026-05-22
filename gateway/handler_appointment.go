package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	pb "hospital/generated/proto"
)

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
