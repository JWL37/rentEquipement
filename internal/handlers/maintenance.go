package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"rentEquipement/internal/repository/maintenance"
	"rentEquipement/internal/session"

	"github.com/go-chi/chi/v5"
)

type MaintenanceHandler struct {
	Repo     maintenance.MaintenanceRepo
	Tmpl     *template.Template
	Sessions *session.SessionsManager
}

func (h *MaintenanceHandler) List(w http.ResponseWriter, r *http.Request) {
	equipmentID := chi.URLParam(r, "equipment_ID")
	if equipmentID == "" {
		http.Error(w, "not id", http.StatusBadRequest)
		return
	}
	maintenances, err := h.Repo.List(equipmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(maintenances); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
