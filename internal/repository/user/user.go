package user

import (
	"rentEquipement/internal/models"
)

type UserRepo interface {
	CreateUser(string, []byte, string, string, string) (*models.User, error)
	Authorize(string) (*models.User, error)
	GetCustomerByUsername(string) (*customerInfo, error)
	GetLogs() ([]models.Log, error)
	GetListOrdersForCustomer(string) ([]OrderView, error)
}

type customerInfo struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Address     string `json:"address,omitempty"`
}
