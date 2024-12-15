package models

import "time"

type Review struct {
	ID          int       `json:"id"`
	CustomerID  int       `json:"customer_id"`
	EquipmentID int       `json:"equipment_id"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	ReviewDate  time.Time `json:"review_date"`
}
