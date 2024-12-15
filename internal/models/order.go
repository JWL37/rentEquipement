package models

import (
	"time"
)

type Order struct {
	ID          int       `json:"id" db:"id"`                     // Уникальный идентификатор заказа
	CustomerID  int       `json:"customer_id" db:"customer_id"`   // Внешний ключ на customer
	EquipmentID int       `json:"equipment_id" db:"equipment_id"` // Внешний ключ на equipment
	StartDate   time.Time `json:"start_date" db:"start_date"`     // Дата начала аренды
	EndDate     time.Time `json:"end_date" db:"end_date"`         // Дата окончания аренды
	TotalCost   float64   `json:"total_cost" db:"total_cost"`     // Общая стоимость
}
