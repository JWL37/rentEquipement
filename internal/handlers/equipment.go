package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"rentEquipement/internal/repository/equipment"
	"rentEquipement/internal/session"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type EquipmentHandler struct {
	Repo     equipment.EquipmentRepo
	Tmpl     *template.Template
	Sessions *session.SessionsManager
}

func (h *EquipmentHandler) FrontListEquipments(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "main.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}
func (h *EquipmentHandler) FrontEquipmentInfo(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "equipment.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *EquipmentHandler) ListEquipments(w http.ResponseWriter, r *http.Request) {
	equipments, err := h.Repo.ListAvailableEquipments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(equipments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *EquipmentHandler) ListNotEquipments(w http.ResponseWriter, r *http.Request) {
	equipments, err := h.Repo.ListNotAvailableEquipments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(equipments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *EquipmentHandler) EquipmentInfo(w http.ResponseWriter, r *http.Request) {
	equipmentID := chi.URLParam(r, "equipment_ID")
	log.Println("equipment_ID: ", equipmentID)
	if equipmentID == "" {
		http.Error(w, "not id", http.StatusBadRequest)
		return
	}
	equipment, err := h.Repo.EquipmentInfo(equipmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(*equipment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *EquipmentHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	log.Printf(sess.ID, sess.Role, sess.UserID)

	if err != nil {
		log.Printf("Create. Error: %s. \n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if sess.Role != "customer" {
		http.Error(w, "not customer", http.StatusForbidden)
		return
	}

	equipmentID := r.FormValue("equipment_ID")
	log.Println(equipmentID)
	startRentStr := r.FormValue("startRent")
	log.Println(startRentStr)
	countDayStr := r.FormValue("countDay")
	log.Println(countDayStr)

	startDate, err := time.Parse("2006-01-02", startRentStr)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	countDays, err := strconv.Atoi(countDayStr)
	if err != nil {
		http.Error(w, "Invalid count of days", http.StatusBadRequest)
		return
	}

	fmt.Println(equipmentID, startDate, countDayStr)

	orderInfo, err := h.Repo.CreateOrder(equipmentID, sess.UserID, startDate, countDays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(*orderInfo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *EquipmentHandler) FormCreateOrder(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		log.Printf("Create. Error: %s. \n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if sess.Role != "customer" {
		http.Error(w, "not customer", http.StatusForbidden)
		return
	}

	err = h.Tmpl.ExecuteTemplate(w, "order.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *EquipmentHandler) DoAvailable(w http.ResponseWriter, r *http.Request) {
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		log.Printf("Create. Error: %s. \n", err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if sess.Role != "admin" {
		http.Error(w, "not admin", http.StatusForbidden)
		return
	}

	equipmentID := chi.URLParam(r, "equipment_ID")
	log.Println("equipment_ID: ", equipmentID)
	if equipmentID == "" {
		http.Error(w, "not id", http.StatusBadRequest)
		return
	}
	err = h.Repo.DoAvailable(equipmentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/info", http.StatusFound)

}
