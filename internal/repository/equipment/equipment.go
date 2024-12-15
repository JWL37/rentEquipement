package equipment

import (
	"rentEquipement/internal/models"
	"time"
)

type EquipmentRepo interface {
	ListAvailableEquipments() ([]models.Equipment, error)
	ListNotAvailableEquipments() ([]models.Equipment, error)
	DoAvailable(string) error
	EquipmentInfo(string) (*models.Equipment, error)
	CreateOrder(string, string, time.Time, int) (*OrderInfo, error)
	// OrderInfo(string)
}

type OrderInfo struct {
	ID        int       `json:"order_id" db:"order_id"`
	StartDate time.Time `json:"start_date" db:"start_date"`
	EndDate   time.Time `json:"end_date" db:"end_date"`
	TotalCost float64   `json:"total_cost" db:"total_cost"`
}
