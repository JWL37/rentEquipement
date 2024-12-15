package models

type Equipment struct {
	ID          int     `json:"id"`          // Уникальный идентификатор оборудования
	Name        string  `json:"name"`        // Название оборудования
	Type        string  `json:"type"`        // Тип оборудования
	Description string  `json:"description"` // Описание оборудования
	PricePerDay float64 `json:"pricePerDay"` // Цена аренды за день
	IsAvailable bool    `json:"isAvailable"` // Статус доступности
}
