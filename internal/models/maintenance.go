package models

import (
	"time"
)

// Maintenance represents a maintenance record in the database.
type Maintenance struct {
	ID          int       `json:"id"`           // Unique identifier of the maintenance
	EquipmentID int       `json:"equipment_id"` // Foreign key to the equipment
	Date        time.Time `json:"date"`         // Date of the maintenance
	Description string    `json:"description"`  // Description of the maintenance work
}
