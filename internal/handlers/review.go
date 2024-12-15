package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"rentEquipement/internal/repository/review"
	"rentEquipement/internal/session"

	"github.com/go-chi/chi/v5"
)

type ReviewHandler struct {
	Repo     review.ReviewRepo
	Tmpl     *template.Template
	Sessions *session.SessionsManager
}

func (h *ReviewHandler) List(w http.ResponseWriter, r *http.Request) {
	equipmentID := chi.URLParam(r, "equipment_ID")
	log.Println("equipment_ID: ", equipmentID)

	reviews, err := h.Repo.List(equipmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(reviews); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
