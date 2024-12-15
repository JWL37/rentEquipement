package models

type Customer struct {
	ID          int    `json:"id"`                     // Уникальный идентификатор клиента
	UserID      int    `json:"user_id"`                // Внешний ключ на пользователя
	Email       string `json:"email"`                  // Email клиента
	PhoneNumber string `json:"phone_number,omitempty"` // Номер телефона клиента
	Address     string `json:"address,omitempty"`      // Адрес клиента
}
