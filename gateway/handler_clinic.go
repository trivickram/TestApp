package main

import (
	"net/http"

	pb "hospital/generated/proto"
)

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
